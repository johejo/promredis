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
	if err := prometheus.Register(promredis.NewCollector(c)); err != nil {
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
	//# HELP go_redis_pool_hits Number of times free connection was found in the pool.
	//# TYPE go_redis_pool_hits gauge
	//go_redis_pool_hits 0
	//# HELP go_redis_pool_idle_conns Number of idle connections in the pool.
	//# TYPE go_redis_pool_idle_conns gauge
	//go_redis_pool_idle_conns 0
	//# HELP go_redis_pool_misses Number of times free connection was NOT found in the pool.
	//# TYPE go_redis_pool_misses gauge
	//go_redis_pool_misses 0
	//# HELP go_redis_pool_stale_conns Number of stale connections removed from the pool.
	//# TYPE go_redis_pool_stale_conns gauge
	//go_redis_pool_stale_conns 0
	//# HELP go_redis_pool_timeouts Number of times a wait timeout occurred.
	//# TYPE go_redis_pool_timeouts gauge
	//go_redis_pool_timeouts 0
	//# HELP go_redis_pool_total_conns Number of total connections in the pool.
	//# TYPE go_redis_pool_total_conns gauge
	//go_redis_pool_total_conns 0
}
