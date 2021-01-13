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
