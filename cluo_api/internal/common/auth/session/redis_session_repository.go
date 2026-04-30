package session

import (
	"context"
	"fmt"
	"time"

	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/google/uuid"
	"github.com/hengadev/errsx"
	"github.com/redis/go-redis/v9"
)

const (
	SessionKeyPrefix          = "auth:session:"
	AccessTokenKeyPrefix      = "auth:access:"
	RefreshTokenKeyPrefix     = "auth:refresh:"
	UserSessionIndexKeyPrefix = "auth:user_sessions:"
)

// RedisSessionRepository implements the minimal SessionRepository interface for authentication middleware
type RedisSessionRepository struct {
	client *redis.Client
}

// NewRedisSessionRepository creates a new Redis-based session repository for authentication
// Returns *RedisSessionRepository which implements both SessionRepository and AuthSessionRepository
func NewRedisSessionRepository(client *redis.Client) *RedisSessionRepository {
	return &RedisSessionRepository{
		client: client,
	}
}

// FindSessionByAccessTokenHash implements two-step security: access token -> session ID -> session data
func (r *RedisSessionRepository) FindSessionByAccessTokenHash(ctx context.Context, accessTokenHash string) (string, []byte, error) {
	// Step 1: Get session ID from access token hash
	accessTokenKey := FormatAccessTokenKey(accessTokenHash)
	sessionID, err := r.client.Get(ctx, accessTokenKey).Result()
	if err != nil {
		return "", nil, errs.ClassifyRedisError("find session ID by access token", err)
	}

	// Step 2: Get session data using session ID
	sessionKey := FormatSessionKey(sessionID)
	sessionData, err := r.client.Get(ctx, sessionKey).Bytes()
	if err != nil {
		return "", nil, errs.ClassifyRedisError("find session data by session ID", err)
	}

	return sessionID, sessionData, nil
}

// FindSessionByRefreshTokenHash implements two-step security: refresh token -> session ID -> session data
func (r *RedisSessionRepository) FindSessionByRefreshTokenHash(ctx context.Context, refreshTokenHash string) (string, []byte, error) {
	// Step 1: Get session ID from refresh token hash
	refreshTokenKey := FormatRefreshTokenKey(refreshTokenHash)
	sessionID, err := r.client.Get(ctx, refreshTokenKey).Result()
	if err != nil {
		return "", nil, errs.ClassifyRedisError("find session ID by refresh token", err)
	}

	// Step 2: Get session data using session ID
	sessionKey := FormatSessionKey(sessionID)
	sessionData, err := r.client.Get(ctx, sessionKey).Result()
	if err != nil {
		return "", nil, errs.ClassifyRedisError("find session data by session ID", err)
	}

	return sessionID, []byte(sessionData), nil
}

// Redis key formatting functions - exported for test helpers
func FormatSessionKey(sessionID string) string {
	return fmt.Sprintf("%s%s", SessionKeyPrefix, sessionID)
}

func FormatAccessTokenKey(accessTokenHash string) string {
	return fmt.Sprintf("%s%s", AccessTokenKeyPrefix, accessTokenHash)
}

func FormatRefreshTokenKey(refreshTokenHash string) string {
	return fmt.Sprintf("%s%s", RefreshTokenKeyPrefix, refreshTokenHash)
}

func FormatUserSessionIndexKey(userIDHash string) string {
	return fmt.Sprintf("%s%s", UserSessionIndexKeyPrefix, userIDHash)
}

// CreateSession creates a new session with access and refresh tokens
func (r *RedisSessionRepository) CreateSession(ctx context.Context, sessionID uuid.UUID, accessTokenHash, refreshTokenHash, userIDHash string, sessionEncoded []byte, accessTTL, refreshTTL time.Duration) error {
	sid := sessionID.String()
	var keysWritten []string

	sessionKey := FormatSessionKey(sid)
	if err := r.client.Set(ctx, sessionKey, sessionEncoded, refreshTTL).Err(); err != nil {
		return errs.ClassifyRedisError("create session data", err)
	}
	keysWritten = append(keysWritten, sessionKey)

	accessTokenKey := FormatAccessTokenKey(accessTokenHash)
	if err := r.client.Set(ctx, accessTokenKey, sid, accessTTL).Err(); err != nil {
		return r.rollbackCreateSession(ctx, keysWritten, errs.ClassifyRedisError("create access token mapping", err))
	}
	keysWritten = append(keysWritten, accessTokenKey)

	refreshTokenKey := FormatRefreshTokenKey(refreshTokenHash)
	if err := r.client.Set(ctx, refreshTokenKey, sid, refreshTTL).Err(); err != nil {
		return r.rollbackCreateSession(ctx, keysWritten, errs.ClassifyRedisError("create refresh token mapping", err))
	}
	keysWritten = append(keysWritten, refreshTokenKey)

	userSessionIndexKey := FormatUserSessionIndexKey(userIDHash)
	if err := r.client.SAdd(ctx, userSessionIndexKey, sid).Err(); err != nil {
		return r.rollbackCreateSession(ctx, keysWritten, errs.ClassifyRedisError("add session to user index", err))
	}

	return nil
}

func (r *RedisSessionRepository) rollbackCreateSession(ctx context.Context, keys []string, originalErr error) error {
	var cleanupErrs errsx.Map
	for _, key := range keys {
		if err := r.client.Del(ctx, key).Err(); err != nil {
			cleanupErrs.Set(key, err)
		}
	}
	if cleanupErrs.Len() > 0 {
		cleanupErrs.Set("original_error", originalErr)
		return cleanupErrs.AsError()
	}
	return originalErr
}

// FindSessionByID retrieves session data by session ID
func (r *RedisSessionRepository) FindSessionByID(ctx context.Context, sessionID uuid.UUID) ([]byte, error) {
	sessionKey := FormatSessionKey(sessionID.String())
	sessionData, err := r.client.Get(ctx, sessionKey).Bytes()
	if err != nil {
		return nil, errs.ClassifyRedisError("find session by ID", err)
	}
	return sessionData, nil
}

// RefreshTokenPair replaces the token pair for a session in Redis.
func (r *RedisSessionRepository) RefreshTokenPair(ctx context.Context, oldRefreshTokenHash, newAccessTokenHash, newRefreshTokenHash string, sessionID uuid.UUID, updatedSessionData []byte, accessTTL, refreshTTL time.Duration) error {
	sid := sessionID.String()

	sessionKey := FormatSessionKey(sid)
	if err := r.client.Set(ctx, sessionKey, updatedSessionData, 0).Err(); err != nil {
		return errs.ClassifyRedisError("update session data", err)
	}

	oldRefreshTokenKey := FormatRefreshTokenKey(oldRefreshTokenHash)
	if err := r.client.Del(ctx, oldRefreshTokenKey).Err(); err != nil {
		return errs.ClassifyRedisError("delete old refresh token", err)
	}

	newAccessTokenKey := FormatAccessTokenKey(newAccessTokenHash)
	if err := r.client.Set(ctx, newAccessTokenKey, sid, accessTTL).Err(); err != nil {
		if delErr := r.client.Del(ctx, oldRefreshTokenKey).Err(); delErr != nil {
			logger, _ := ctxutil.GetLoggerFromContext(ctx)
			if logger != nil {
				logger.WarnContext(ctx, "failed to restore old refresh token during rollback", "error", delErr)
			}
		}
		return errs.ClassifyRedisError("create new access token mapping", err)
	}

	newRefreshTokenKey := FormatRefreshTokenKey(newRefreshTokenHash)
	if err := r.client.Set(ctx, newRefreshTokenKey, sid, refreshTTL).Err(); err != nil {
		if delErr := r.client.Del(ctx, newAccessTokenKey).Err(); delErr != nil {
			logger, _ := ctxutil.GetLoggerFromContext(ctx)
			if logger != nil {
				logger.WarnContext(ctx, "failed to cleanup new access token during rollback", "error", delErr)
			}
		}
		if delErr := r.client.Del(ctx, oldRefreshTokenKey).Err(); delErr != nil {
			logger, _ := ctxutil.GetLoggerFromContext(ctx)
			if logger != nil {
				logger.WarnContext(ctx, "failed to restore old refresh token during rollback", "error", delErr)
			}
		}
		return errs.ClassifyRedisError("create new refresh token mapping", err)
	}

	r.cleanupStaleAccessTokens(ctx, sid, newAccessTokenKey)

	return nil
}

func (r *RedisSessionRepository) cleanupStaleAccessTokens(ctx context.Context, sessionID, newAccessTokenKey string) {
	logger, _ := ctxutil.GetLoggerFromContext(ctx)

	iter := r.client.Scan(ctx, 0, AccessTokenKeyPrefix+"*", 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		if key == newAccessTokenKey {
			continue
		}
		val, err := r.client.Get(ctx, key).Result()
		if err != nil {
			if logger != nil {
				logger.WarnContext(ctx, "failed to check access token for cleanup", "key", key, "error", err)
			}
			continue
		}
		if val == sessionID {
			if delErr := r.client.Del(ctx, key).Err(); delErr != nil && logger != nil {
				logger.WarnContext(ctx, "failed to delete stale access token", "key", key, "error", delErr)
			}
		}
	}
	if err := iter.Err(); err != nil && logger != nil {
		logger.WarnContext(ctx, "error iterating access tokens for cleanup", "error", err)
	}
}

// RemoveSessionByID removes a session and its token mappings by session ID.
func (r *RedisSessionRepository) RemoveSessionByID(ctx context.Context, sessionID uuid.UUID) error {
	sid := sessionID.String()
	sessionKey := FormatSessionKey(sid)

	sessionData, err := r.client.Get(ctx, sessionKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil
		}
		return errs.ClassifyRedisError("get session data for removal", err)
	}

	sessionEncx, err := DecodeSession([]byte(sessionData))
	if err != nil {
		return errs.ClassifyRedisError("decode session for removal", err)
	}

	pipe := r.client.Pipeline()
	pipe.Del(ctx, sessionKey)
	pipe.Del(ctx, FormatAccessTokenKey(sessionEncx.AccessTokenHash))
	pipe.Del(ctx, FormatRefreshTokenKey(sessionEncx.RefreshTokenHash))
	pipe.SRem(ctx, FormatUserSessionIndexKey(sessionEncx.UserIDHash), sid)

	if _, err := pipe.Exec(ctx); err != nil {
		return errs.ClassifyRedisError("remove session", err)
	}

	return nil
}
