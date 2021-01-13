// Package promredis exports pool stats of go-redis/redis as prometheus metrics.
package promredis

import (
	"github.com/go-redis/redis/v8"
	"github.com/prometheus/client_golang/prometheus"
)

// Client describes a getter for *redis.PoolStats.
type Client interface {
	PoolStats() *redis.PoolStats
}

var _ Client = (*redis.Client)(nil)

type collector struct {
	client Client

	hitsDesc       *prometheus.Desc
	missesDesc     *prometheus.Desc
	timeoutsDesc   *prometheus.Desc
	totalConnsDesc *prometheus.Desc
	idleConnsDesc  *prometheus.Desc
	staleConnsDesc *prometheus.Desc
}

var _ prometheus.Collector = (*collector)(nil)

// NewCollector returns a new collector implements prometheus.Collector.
func NewCollector(client Client) prometheus.Collector {
	return &collector{
		client:         client,
		hitsDesc:       prometheus.NewDesc(fqName("hits"), "Number of times free connection was found in the pool.", nil, nil),
		missesDesc:     prometheus.NewDesc(fqName("misses"), "Number of times free connection was NOT found in the pool.", nil, nil),
		timeoutsDesc:   prometheus.NewDesc(fqName("timeouts"), "Number of times a wait timeout occurred.", nil, nil),
		totalConnsDesc: prometheus.NewDesc(fqName("total_conns"), "Number of total connections in the pool.", nil, nil),
		idleConnsDesc:  prometheus.NewDesc(fqName("idle_conns"), "Number of idle connections in the pool.", nil, nil),
		staleConnsDesc: prometheus.NewDesc(fqName("stale_conns"), "Number of stale connections removed from the pool.", nil, nil),
	}
}

func fqName(name string) string {
	return "go_redis_pool_stats_" + name
}

// Describe implements prometheus.Collector.
func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.hitsDesc
	ch <- c.missesDesc
	ch <- c.timeoutsDesc
	ch <- c.totalConnsDesc
	ch <- c.idleConnsDesc
	ch <- c.staleConnsDesc
}

// Collect implements prometheus.Collector.
func (c *collector) Collect(ch chan<- prometheus.Metric) {
	stats := c.client.PoolStats()
	ch <- prometheus.MustNewConstMetric(c.hitsDesc, prometheus.GaugeValue, float64(stats.Hits))
	ch <- prometheus.MustNewConstMetric(c.missesDesc, prometheus.GaugeValue, float64(stats.Misses))
	ch <- prometheus.MustNewConstMetric(c.timeoutsDesc, prometheus.GaugeValue, float64(stats.Timeouts))
	ch <- prometheus.MustNewConstMetric(c.totalConnsDesc, prometheus.GaugeValue, float64(stats.TotalConns))
	ch <- prometheus.MustNewConstMetric(c.idleConnsDesc, prometheus.GaugeValue, float64(stats.IdleConns))
	ch <- prometheus.MustNewConstMetric(c.staleConnsDesc, prometheus.GaugeValue, float64(stats.StaleConns))
}
