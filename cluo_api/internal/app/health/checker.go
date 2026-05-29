package health

import (
	"context"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// Status represents the health status of a component.
type Status string

const (
	StatusHealthy   Status = "healthy"
	StatusUnhealthy Status = "unhealthy"
	StatusDegraded  Status = "degraded"
)

// CheckResult represents the result of a health check.
type CheckResult struct {
	Status      Status `json:"status"`
	Message     string `json:"message,omitempty"`
	DurationMs  int64  `json:"duration_ms"`
}

// HealthReport represents the overall health of the system.
type HealthReport struct {
	Status        Status                 `json:"status"`
	Checks        map[string]CheckResult `json:"checks"`
	Version       string                 `json:"version,omitempty"`
	Timestamp     time.Time              `json:"timestamp"`
	TotalDurationMs int64                `json:"total_duration_ms"`
}

// Checker performs health checks on system components.
type Checker struct {
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	version     string
	timeout     time.Duration
}

// NewChecker creates a new health checker.
func NewChecker(dbPool *pgxpool.Pool, redisClient *redis.Client, version string) *Checker {
	return &Checker{
		dbPool:      dbPool,
		redisClient: redisClient,
		version:     version,
		timeout:     5 * time.Second,
	}
}

// CheckLiveness performs a basic liveness check.
// This should always return healthy if the server is running.
func (c *Checker) CheckLiveness() HealthReport {
	return HealthReport{
		Status:    StatusHealthy,
		Checks:    map[string]CheckResult{},
		Version:   c.version,
		Timestamp: time.Now().UTC(),
	}
}

// CheckReadiness performs readiness checks on all dependencies.
func (c *Checker) CheckReadiness(ctx context.Context) HealthReport {
	start := time.Now()

	report := HealthReport{
		Status:    StatusHealthy,
		Checks:    make(map[string]CheckResult),
		Version:   c.version,
		Timestamp: time.Now().UTC(),
	}

	// Create a context with timeout
	checkCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	// Run checks concurrently
	var wg sync.WaitGroup
	var mu sync.Mutex

	// Database check
	if c.dbPool != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result := c.checkDatabase(checkCtx)
			mu.Lock()
			report.Checks["database"] = result
			mu.Unlock()
		}()
	}

	// Redis check
	if c.redisClient != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result := c.checkRedis(checkCtx)
			mu.Lock()
			report.Checks["redis"] = result
			mu.Unlock()
		}()
	}

	wg.Wait()

	// Determine overall status
	report.Status = c.calculateOverallStatus(report.Checks)
	report.TotalDurationMs = time.Since(start).Milliseconds()

	return report
}

func (c *Checker) checkDatabase(ctx context.Context) CheckResult {
	start := time.Now()

	err := c.dbPool.Ping(ctx)
	durationMs := time.Since(start).Milliseconds()

	if err != nil {
		return CheckResult{
			Status:     StatusUnhealthy,
			Message:    err.Error(),
			DurationMs: durationMs,
		}
	}

	// Check pool stats for degraded state
	stats := c.dbPool.Stat()
	if stats.AcquiredConns() > int32(float64(stats.MaxConns())*0.9) {
		return CheckResult{
			Status:     StatusDegraded,
			Message:    "connection pool near capacity",
			DurationMs: durationMs,
		}
	}

	return CheckResult{
		Status:     StatusHealthy,
		DurationMs: durationMs,
	}
}

func (c *Checker) checkRedis(ctx context.Context) CheckResult {
	start := time.Now()

	err := c.redisClient.Ping(ctx).Err()
	durationMs := time.Since(start).Milliseconds()

	if err != nil {
		return CheckResult{
			Status:     StatusUnhealthy,
			Message:    err.Error(),
			DurationMs: durationMs,
		}
	}

	return CheckResult{
		Status:     StatusHealthy,
		DurationMs: durationMs,
	}
}

func (c *Checker) calculateOverallStatus(checks map[string]CheckResult) Status {
	hasUnhealthy := false
	hasDegraded := false

	for _, check := range checks {
		switch check.Status {
		case StatusUnhealthy:
			hasUnhealthy = true
		case StatusDegraded:
			hasDegraded = true
		}
	}

	if hasUnhealthy {
		return StatusUnhealthy
	}
	if hasDegraded {
		return StatusDegraded
	}
	return StatusHealthy
}
