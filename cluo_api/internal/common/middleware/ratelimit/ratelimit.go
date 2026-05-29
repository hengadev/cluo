// Package ratelimit provides an IP-based rate-limiting middleware using a sliding
// window algorithm backed by an in-memory store. The store is defined behind a
// port interface (RateLimiterStore) so it can be swapped to Redis later.
package ratelimit

import (
	"context"
	"fmt"
	"math"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

// ---------------------------------------------------------------------------
// Port interface
// ---------------------------------------------------------------------------

// RateLimiterStore is the port that any backing store must satisfy.
// Implementations must be safe for concurrent use.
type RateLimiterStore interface {
	// Allow checks whether a request from key is allowed under the given limit
	// (maxRequests in window). It returns (allowed bool, retryAfter time.Duration).
	// retryAfter is meaningful only when allowed is false.
	Allow(key string, maxRequests int, window time.Duration) (bool, time.Duration)
}

// ---------------------------------------------------------------------------
// In-memory implementation
// ---------------------------------------------------------------------------

// entry tracks a single key's request timestamps within the sliding window.
type entry struct {
	timestamps []time.Time
}

// InMemoryStore is a RateLimiterStore backed by an in-memory map.
type InMemoryStore struct {
	mu   sync.Mutex
	data map[string]*entry
}

// NewInMemoryStore creates a new InMemoryStore.
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{data: make(map[string]*entry)}
}

// Allow implements RateLimiterStore.
func (s *InMemoryStore) Allow(key string, maxRequests int, window time.Duration) (bool, time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-window)

	e, ok := s.data[key]
	if !ok {
		e = &entry{}
		s.data[key] = e
	}

	// Filter out timestamps outside the window (sliding window).
	valid := e.timestamps[:0]
	for _, t := range e.timestamps {
		if t.After(windowStart) {
			valid = append(valid, t)
		}
	}
	e.timestamps = valid

	if len(e.timestamps) >= maxRequests {
		// The oldest timestamp in the window determines when the window resets.
		oldest := e.timestamps[0]
		retryAfter := oldest.Add(window).Sub(now)
		if retryAfter < 0 {
			retryAfter = 0
		}
		return false, retryAfter
	}

	e.timestamps = append(e.timestamps, now)
	return true, 0
}

// Cleanup removes expired entries. Call periodically to prevent unbounded memory growth.
func (s *InMemoryStore) Cleanup(window time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-window)

	for key, e := range s.data {
		// Remove old timestamps.
		valid := e.timestamps[:0]
		for _, t := range e.timestamps {
			if t.After(windowStart) {
				valid = append(valid, t)
			}
		}
		if len(valid) == 0 {
			delete(s.data, key)
		} else {
			e.timestamps = valid
		}
	}
}

// StartCleanup starts a background goroutine that calls Cleanup on the given
// interval. It stops when ctx is cancelled.
func (s *InMemoryStore) StartCleanup(ctx context.Context, window time.Duration) {
	go func() {
		ticker := time.NewTicker(window)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				s.Cleanup(window)
			}
		}
	}()
}

// ---------------------------------------------------------------------------
// Middleware
// ---------------------------------------------------------------------------

// Config holds the rate-limiting configuration for a named route group.
type Config struct {
	// MaxRequests is the maximum number of requests allowed in the window.
	MaxRequests int
	// Window is the time window duration.
	Window time.Duration
}

// RateLimiter returns middleware that rate-limits requests per client IP
// using the provided store and config.
func RateLimiter(store RateLimiterStore, cfg Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := clientIP(r)

			allowed, retryAfter := store.Allow(ip, cfg.MaxRequests, cfg.Window)
			if !allowed {
				w.Header().Set("Retry-After", fmt.Sprintf("%d", int(math.Ceil(retryAfter.Seconds()))))
				http.Error(w, "too many requests", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// clientIP extracts the client IP from the request.
// It prefers X-Forwarded-For (first value, since traffic passes through Caddy)
// with a fallback to RemoteAddr.
func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		if ip := strings.TrimSpace(strings.Split(xff, ",")[0]); ip != "" {
			if net.ParseIP(ip) != nil {
				return ip
			}
		}
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
