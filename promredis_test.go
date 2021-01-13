package promredis

import (
	"bufio"
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func TestRegister(t *testing.T) {
	ctx := context.Background()

	rc := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer rc.Close()

	if err := rc.Ping(ctx).Err(); err != nil {
		t.Fatal(err)
	}
	c := NewCollector(rc)
	if err := prometheus.Register(c); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { prometheus.Unregister(c) })

	for i := 0; i < 100; i++ {
		k := "key" + strconv.Itoa(i)
		rc.Set(ctx, k, i, -1)
		rc.Get(ctx, k)
	}

	time.Sleep(1 * time.Second)

	ph := promhttp.Handler()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	ph.ServeHTTP(rec, req)

	resp := rec.Result()

	hitsName := fqName("hits")
	missesName := fqName("misses")
	timeoutsName := fqName("timeouts")
	totalConnsName := fqName("total_conns")
	idleConnsName := fqName("idle_conns")
	staleConnsName := fqName("stale_conns")

	all := []string{
		hitsName,
		missesName,
		timeoutsName,
		totalConnsName,
		idleConnsName,
		staleConnsName,
	}
	type result struct {
		found bool
		value float64
	}
	results := make(map[string]result)
	for _, name := range all {
		results[name] = result{found: false}
	}

	s := bufio.NewScanner(resp.Body)
	for s.Scan() {
		line := s.Text()
		println(line)
		if strings.HasPrefix(line, "#") {
			continue
		}
		for _, name := range all {
			if strings.Contains(line, name) {
				value, err := strconv.ParseFloat(strings.Split(line, " ")[1], 10)
				if err != nil {
					t.Fatal(err)
				}
				results[name] = result{found: true, value: value}
			}
		}
	}

	for name, result := range results {
		if !result.found {
			t.Errorf("%s is not found", name)
		}
	}

	if v := results[hitsName].value; v == 0 {
		t.Errorf("%s should not be 0, but got %v", hitsName, v)
	}
}
