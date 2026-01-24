package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics holds all application metrics.
type Metrics struct {
	// HTTP metrics
	HTTPRequestsTotal    *prometheus.CounterVec
	HTTPRequestDuration  *prometheus.HistogramVec
	HTTPRequestsInFlight prometheus.Gauge

	// Database metrics
	DBConnectionsOpen prometheus.GaugeFunc
	DBConnectionsIdle prometheus.GaugeFunc

	// Business metrics
	CasesCreated   prometheus.Counter
	ClientsCreated prometheus.Counter
}

// Config holds metrics configuration.
type Config struct {
	Namespace   string
	Subsystem   string
	ConstLabels prometheus.Labels
}

// New creates a new metrics instance with all metrics registered.
func New(cfg Config) *Metrics {
	if cfg.Namespace == "" {
		cfg.Namespace = "cluo"
	}
	if cfg.Subsystem == "" {
		cfg.Subsystem = "api"
	}

	m := &Metrics{}

	// HTTP request counter
	m.HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace:   cfg.Namespace,
			Subsystem:   cfg.Subsystem,
			Name:        "http_requests_total",
			Help:        "Total number of HTTP requests",
			ConstLabels: cfg.ConstLabels,
		},
		[]string{"method", "path", "status"},
	)

	// HTTP request duration histogram
	m.HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace:   cfg.Namespace,
			Subsystem:   cfg.Subsystem,
			Name:        "http_request_duration_seconds",
			Help:        "HTTP request duration in seconds",
			ConstLabels: cfg.ConstLabels,
			Buckets:     []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
		},
		[]string{"method", "path"},
	)

	// HTTP requests in flight gauge
	m.HTTPRequestsInFlight = promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace:   cfg.Namespace,
			Subsystem:   cfg.Subsystem,
			Name:        "http_requests_in_flight",
			Help:        "Number of HTTP requests currently being processed",
			ConstLabels: cfg.ConstLabels,
		},
	)

	// Business metrics
	m.CasesCreated = promauto.NewCounter(
		prometheus.CounterOpts{
			Namespace:   cfg.Namespace,
			Subsystem:   cfg.Subsystem,
			Name:        "cases_created_total",
			Help:        "Total number of cases created",
			ConstLabels: cfg.ConstLabels,
		},
	)

	m.ClientsCreated = promauto.NewCounter(
		prometheus.CounterOpts{
			Namespace:   cfg.Namespace,
			Subsystem:   cfg.Subsystem,
			Name:        "clients_created_total",
			Help:        "Total number of clients created",
			ConstLabels: cfg.ConstLabels,
		},
	)

	return m
}

// DBPoolStats provides database pool statistics.
type DBPoolStats interface {
	TotalConns() int32
	IdleConns() int32
	AcquiredConns() int32
	MaxConns() int32
}

// RegisterDBMetrics registers database pool metrics.
func (m *Metrics) RegisterDBMetrics(namespace, subsystem string, stats func() DBPoolStats) {
	if stats == nil {
		return
	}

	// Database connections open
	m.DBConnectionsOpen = prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "db_connections_open",
			Help:      "Number of open database connections",
		},
		func() float64 {
			return float64(stats().TotalConns())
		},
	)
	prometheus.MustRegister(m.DBConnectionsOpen)

	// Database connections idle
	m.DBConnectionsIdle = prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "db_connections_idle",
			Help:      "Number of idle database connections",
		},
		func() float64 {
			return float64(stats().IdleConns())
		},
	)
	prometheus.MustRegister(m.DBConnectionsIdle)
}
