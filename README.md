# promredis

**DEPRECATED** Use https://pkg.go.dev/github.com/go-redis/redis/extra/redisprometheus/v9 instead

[![ci](https://github.com/johejo/promredis/workflows/ci/badge.svg?branch=main)](https://github.com/johejo/promredis/actions?query=workflow%3Aci)
[![Go Reference](https://pkg.go.dev/badge/github.com/johejo/promredis.svg)](https://pkg.go.dev/github.com/johejo/promredis)
[![codecov](https://codecov.io/gh/johejo/promredis/branch/main/graph/badge.svg)](https://codecov.io/gh/johejo/promredis)
[![Go Report Card](https://goreportcard.com/badge/github.com/johejo/promredis)](https://goreportcard.com/report/github.com/johejo/promredis)

Package promredis exports pool stats of go-redis/redis as prometheus metrics.

## Example

```go
package promredis_test

import (
	"bufio"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/johejo/promredis"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Example() {
	c := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	if err := prometheus.Register(promredis.NewPoolStatsCollector(c)); err != nil {
		panic(err)
	}

	go func() {
		http.ListenAndServe(":8080", promhttp.Handler())
	}()

	time.Sleep(1 * time.Second)

	resp, err := http.Get("http://localhost:8080")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	s := bufio.NewScanner(resp.Body)

	for s.Scan() {
		line := s.Text()
		if strings.Contains(line, "redis") {
			fmt.Println(line)
		}
	}

	// Output:
	//# HELP go_redis_pool_stats_hits Number of times free connection was found in the pool.
	//# TYPE go_redis_pool_stats_hits gauge
	//go_redis_pool_stats_hits 0
	//# HELP go_redis_pool_stats_idle_conns Number of idle connections in the pool.
	//# TYPE go_redis_pool_stats_idle_conns gauge
	//go_redis_pool_stats_idle_conns 0
	//# HELP go_redis_pool_stats_misses Number of times free connection was NOT found in the pool.
	//# TYPE go_redis_pool_stats_misses gauge
	//go_redis_pool_stats_misses 0
	//# HELP go_redis_pool_stats_stale_conns Number of stale connections removed from the pool.
	//# TYPE go_redis_pool_stats_stale_conns gauge
	//go_redis_pool_stats_stale_conns 0
	//# HELP go_redis_pool_stats_timeouts Number of times a wait timeout occurred.
	//# TYPE go_redis_pool_stats_timeouts gauge
	//go_redis_pool_stats_timeouts 0
	//# HELP go_redis_pool_stats_total_conns Number of total connections in the pool.
	//# TYPE go_redis_pool_stats_total_conns gauge
	//go_redis_pool_stats_total_conns 0
}
```

## License

MIT

## Author

Mitsuo Heijo
