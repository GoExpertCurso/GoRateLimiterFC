package main

import (
	"net/http"

	"github.com/GoExpertCurso/GoRateLimiterFC/configs"
	web "github.com/GoExpertCurso/GoRateLimiterFC/internal/infra/web"
	mid "github.com/GoExpertCurso/GoRateLimiterFC/internal/infra/web/middleware"
	"github.com/go-redis/redis"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     configs.DBHost + ":" + configs.DBPort,
		Password: configs.DBPassword,
		DB:       0,
	})
	defer client.Close()

	_, err = client.Ping().Result()
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", web.Home)

	wrappedMux := mid.RateLimitMiddleware(mux, client)

	http.ListenAndServe(":8080", wrappedMux)
}
