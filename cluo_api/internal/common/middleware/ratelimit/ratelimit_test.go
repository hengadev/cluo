package ratelimit

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// ---------------------------------------------------------------------------
// InMemoryStore unit tests
// ---------------------------------------------------------------------------

func TestInMemoryStore_Allow_UnderLimit(t *testing.T) {
	store := NewInMemoryStore()
	allowed, retryAfter := store.Allow("1.2.3.4", 3, time.Minute)
	assert.True(t, allowed)
	assert.Equal(t, time.Duration(0), retryAfter)
}

func TestInMemoryStore_Allow_AtLimit(t *testing.T) {
	store := NewInMemoryStore()
	window := time.Minute
	maxReqs := 2

	store.Allow("1.2.3.4", maxReqs, window)
	store.Allow("1.2.3.4", maxReqs, window)

	allowed, _ := store.Allow("1.2.3.4", maxReqs, window)
	assert.False(t, allowed)
}

func TestInMemoryStore_Allow_DifferentKeysTrackedIndependently(t *testing.T) {
	store := NewInMemoryStore()
	window := time.Minute
	maxReqs := 1

	allowed1, _ := store.Allow("1.1.1.1", maxReqs, window)
	assert.True(t, allowed1)

	allowed2, _ := store.Allow("2.2.2.2", maxReqs, window)
	assert.True(t, allowed2)

	// First key should now be blocked.
	allowed1Again, _ := store.Allow("1.1.1.1", maxReqs, window)
	assert.False(t, allowed1Again)

	// Second key has also reached its own limit (tracked independently from first key).
	allowed2Again, _ := store.Allow("2.2.2.2", maxReqs, window)
	assert.False(t, allowed2Again)
}

func TestInMemoryStore_WindowReset(t *testing.T) {
	store := NewInMemoryStore()
	maxReqs := 1
	window := 100 * time.Millisecond

	store.Allow("1.2.3.4", maxReqs, window)

	// Immediately after, should be blocked.
	allowed, _ := store.Allow("1.2.3.4", maxReqs, window)
	assert.False(t, allowed)

	// Wait for the window to pass.
	time.Sleep(150 * time.Millisecond)

	allowed, _ = store.Allow("1.2.3.4", maxReqs, window)
	assert.True(t, allowed)
}

func TestInMemoryStore_RetryAfter(t *testing.T) {
	store := NewInMemoryStore()
	maxReqs := 1
	window := 10 * time.Second

	store.Allow("1.2.3.4", maxReqs, window)

	_, retryAfter := store.Allow("1.2.3.4", maxReqs, window)
	assert.Greater(t, retryAfter, time.Duration(0))
	assert.LessOrEqual(t, retryAfter, window)
}

func TestInMemoryStore_Cleanup(t *testing.T) {
	store := NewInMemoryStore()
	window := 100 * time.Millisecond

	store.Allow("1.2.3.4", 10, window)

	// Wait for entries to expire.
	time.Sleep(150 * time.Millisecond)

	store.Cleanup(window)

	store.mu.Lock()
	_, exists := store.data["1.2.3.4"]
	store.mu.Unlock()
	assert.False(t, exists, "cleanup should remove expired entries")
}

// ---------------------------------------------------------------------------
// Middleware integration tests
// ---------------------------------------------------------------------------

func TestMiddleware_FirstRequestAllowed(t *testing.T) {
	store := NewInMemoryStore()
	cfg := Config{MaxRequests: 5, Window: time.Minute}

	called := false
	handler := RateLimiter(store, cfg)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "1.2.3.4:1234"
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.True(t, called)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestMiddleware_LimitReached_Returns429(t *testing.T) {
	store := NewInMemoryStore()
	cfg := Config{MaxRequests: 2, Window: time.Minute}

	okHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := RateLimiter(store, cfg)(okHandler)

	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "1.2.3.4:1234"
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	}

	// Third request should be rate limited.
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "1.2.3.4:1234"
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusTooManyRequests, rec.Code)
	assert.NotEmpty(t, rec.Header().Get("Retry-After"))
}

func TestMiddleware_WindowReset_AllowsNewRequests(t *testing.T) {
	store := NewInMemoryStore()
	cfg := Config{MaxRequests: 1, Window: 100 * time.Millisecond}

	okHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := RateLimiter(store, cfg)(okHandler)

	// First request allowed.
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "1.2.3.4:1234"
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Second request blocked.
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "1.2.3.4:1234"
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusTooManyRequests, rec.Code)

	// Wait for window to reset.
	time.Sleep(150 * time.Millisecond)

	// New request should be allowed.
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "1.2.3.4:1234"
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestMiddleware_DifferentIPsTrackedIndependently(t *testing.T) {
	store := NewInMemoryStore()
	cfg := Config{MaxRequests: 1, Window: time.Minute}

	okHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := RateLimiter(store, cfg)(okHandler)

	// First IP: use up the limit.
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "1.1.1.1:1234"
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	// First IP: now blocked.
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "1.1.1.1:1234"
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusTooManyRequests, rec.Code)

	// Second IP: should still be allowed.
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "2.2.2.2:5678"
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestMiddleware_XForwardedFor(t *testing.T) {
	store := NewInMemoryStore()
	cfg := Config{MaxRequests: 1, Window: time.Minute}

	okHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := RateLimiter(store, cfg)(okHandler)

	// Request with X-Forwarded-For header.
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "10.0.0.1:1234" // internal proxy
	req.Header.Set("X-Forwarded-For", "9.8.7.6, 10.0.0.1")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Same XFF IP but different RemoteAddr — should be blocked (same key: 9.8.7.6).
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "10.0.0.2:1234"
	req.Header.Set("X-Forwarded-For", "9.8.7.6, 10.0.0.1")
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusTooManyRequests, rec.Code)
}

// ---------------------------------------------------------------------------
// clientIP tests
// ---------------------------------------------------------------------------

func TestClientIP_XForwardedFor(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "10.0.0.1:1234"
	req.Header.Set("X-Forwarded-For", "1.2.3.4, 5.6.7.8")

	assert.Equal(t, "1.2.3.4", clientIP(req))
}

func TestClientIP_FallbackToRemoteAddr(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "192.168.1.1:4321"

	assert.Equal(t, "192.168.1.1", clientIP(req))
}

func TestClientIP_InvalidXFF_FallbackToRemoteAddr(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "192.168.1.1:4321"
	req.Header.Set("X-Forwarded-For", "not-an-ip")

	assert.Equal(t, "192.168.1.1", clientIP(req))
}
