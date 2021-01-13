// Package promredis exports pool stats of go-redis/redis as prometheus metrics.
package promredis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	hitsName       = "go_redis_pool_stats_hits"
	missesName     = "go_redis_pool_stats_misses"
	timeoutsName   = "go_redis_pool_stats_timeouts"
	totalConnsName = "go_redis_pool_stats_total_conns"
	idleConnsName  = "go_redis_pool_stats_idle_conns"
	staleConnsName = "go_redis_pool_stats_stale_conns"
)

// Client describes a getter for *redis.PoolStats.
type Client interface {
	PoolStats() *redis.PoolStats
}

var _ Client = (*redis.Client)(nil)

func newGauge(namespace, subsystem, name, help string) prometheus.Gauge {
	return prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      name,
		Help:      help,
	})
}

// Register registers prometheus metrics by prometheus.Gauge and returns a finalize func.
// The finalize func unregisters collectors and cleans up resources.
func Register(ctx context.Context, client Client, opts ...Option) (func(), error) {
	cfg := new(config)
	for _, opt := range defaults() {
		opt(cfg)
	}

	hits := newGauge(cfg.namespce, cfg.subsystem, hitsName, "Number of times free connection was found in the pool.")
	misses := newGauge(cfg.namespce, cfg.subsystem, missesName, "Number of times free connection was NOT found in the pool.")
	timeouts := newGauge(cfg.namespce, cfg.subsystem, timeoutsName, "Number of times a wait timeout occurred.")
	totalConns := newGauge(cfg.namespce, cfg.subsystem, totalConnsName, "Number of total connections in the pool.")
	idleConns := newGauge(cfg.namespce, cfg.subsystem, idleConnsName, "Number of idle connections in the pool.")
	staleConns := newGauge(cfg.namespce, cfg.subsystem, staleConnsName, "Number of stale connections removed from the pool.")

	cs := []prometheus.Collector{hits, misses, timeouts, totalConns, idleConns, staleConns}

	for _, c := range cs {
		if err := cfg.registerer.Register(c); err != nil {
			return nil, err
		}
	}
	unregister := func() {
		for _, c := range cs {
			cfg.registerer.Unregister(c)
		}
	}

	ctx, cancel := context.WithCancel(ctx)

	go func() {
		defer cfg.ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-cfg.ticker.Tick():
				stats := client.PoolStats()
				hits.Set(float64(stats.Hits))
				misses.Set(float64(stats.Misses))
				timeouts.Set(float64(stats.Timeouts))
				totalConns.Set(float64(stats.TotalConns))
				idleConns.Set(float64(stats.IdleConns))
				staleConns.Set(float64(stats.StaleConns))
			}
		}
	}()

	return func() {
		unregister()
		cancel()
	}, nil
}

func defaults() []Option {
	return []Option{
		WithDefaultTicker(1 * time.Second),
		WithRegisterer(prometheus.DefaultRegisterer),
	}
}

// Option describes an option for Register.
type Option func(*config)

type config struct {
	namespce, subsystem string
	registerer          prometheus.Registerer
	ticker              Ticker
}

// WithNamespace returns an option sets namespace.
func WithNamespace(namespace string) Option {
	return func(cfg *config) {
		cfg.namespce = namespace
	}
}

// WithSubsystem returns an option sets subsystem.
func WithSubsystem(subsystem string) Option {
	return func(cfg *config) {
		cfg.subsystem = subsystem
	}
}

// WithDefaultTicker return an option that uses time.Ticker.
func WithDefaultTicker(interval time.Duration) Option {
	return WithTicker(&defaultTicker{ticker: time.NewTicker(interval)})
}

// WithDefaultTicker return an option that uses specified Ticker.
func WithTicker(ticker Ticker) Option {
	return func(cfg *config) {
		cfg.ticker = ticker
	}
}

// WithRegisterer returns an option that uses prometheus.Registerer.
func WithRegisterer(registerer prometheus.Registerer) Option {
	return func(cfg *config) {
		cfg.registerer = registerer
	}
}

// Ticker is interface that wraps like time.Ticker.
type Ticker interface {
	// Tick returns an channel of time.Time like time.Ticker.C
	Tick() <-chan time.Time
	// Stop stops ticker like time.Ticker#Stop
	Stop()
}

type defaultTicker struct {
	ticker *time.Ticker
}

var _ Ticker = (*defaultTicker)(nil)

func (t *defaultTicker) Tick() <-chan time.Time {
	return t.ticker.C
}

func (t *defaultTicker) Stop() {
	t.ticker.Stop()
}
