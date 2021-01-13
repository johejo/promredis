# go-redis-prometheus

[![ci](https://github.com/johejo/go-redis-prometheus/workflows/ci/badge.svg?branch=main)](https://github.com/johejo/go-redis-prometheus/actions?query=workflow%3Aci)
[![Go Reference](https://pkg.go.dev/badge/github.com/johejo/go-redis-prometheus.svg)](https://pkg.go.dev/github.com/johejo/go-redis-prometheus)
[![codecov](https://codecov.io/gh/johejo/go-redis-prometheus/branch/main/graph/badge.svg)](https://codecov.io/gh/johejo/go-redis-prometheus)
[![Go Report Card](https://goreportcard.com/badge/github.com/johejo/go-redis-prometheus)](https://goreportcard.com/report/github.com/johejo/go-redis-prometheus)

Package promredis exports pool stats of go-redis/redis as prometheus metrics.

## Example

```go
package promredis_test

import (
	"context"
	"io"
	"net/http"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/johejo/go-redis-prometheus/promredis"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Example() {
	c := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	ctx := context.Background()
	finish, err := promredis.Register(ctx, c)
	if err != nil {
		panic(err)
	}
	defer finish()

	go func() {
		http.ListenAndServe(":8080", promhttp.Handler())
	}()

	resp, err := http.Get("http://localhost:8080")
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, resp.Body)
}
```

Result

```
...
# HELP go_redis_pool_stats_hits Number of times free connection was found in the pool.
# TYPE go_redis_pool_stats_hits gauge
go_redis_pool_stats_hits 200
# HELP go_redis_pool_stats_idle_conns Number of idle connections in the pool.
# TYPE go_redis_pool_stats_idle_conns gauge
go_redis_pool_stats_idle_conns 1
# HELP go_redis_pool_stats_misses Number of times free connection was NOT found in the pool.
# TYPE go_redis_pool_stats_misses gauge
go_redis_pool_stats_misses 1
# HELP go_redis_pool_stats_stale_conns Number of stale connections removed from the pool.
# TYPE go_redis_pool_stats_stale_conns gauge
go_redis_pool_stats_stale_conns 0
# HELP go_redis_pool_stats_timeouts Number of times a wait timeout occurred.
# TYPE go_redis_pool_stats_timeouts gauge
go_redis_pool_stats_timeouts 0
# HELP go_redis_pool_stats_total_conns Number of total connections in the pool.
# TYPE go_redis_pool_stats_total_conns gauge
go_redis_pool_stats_total_conns 1
...
```

## License

MIT

## Author

Mitsuo Heijo
