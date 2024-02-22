package main

import (
	"net/http"

	"github.com/GoExpertCurso/GoRateLimiterFC/configs"
	web "github.com/GoExpertCurso/GoRateLimiterFC/internal/infra/web"
	mid "github.com/GoExpertCurso/GoRateLimiterFC/internal/infra/web/middleware"
	"github.com/go-redis/redis"
)

func main() {
	config, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	conf := configs.NewConf(config.TokenLimit, config.IPLimit)

	client := redis.NewClient(&redis.Options{
		Addr:     config.DBHost + ":" + config.DBPort,
		Password: config.DBPassword,
		DB:       0,
	})
	defer client.Close()

	_, err = client.Ping().Result()
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", web.Home)

	wrappedMux := mid.RateLimitMiddleware(mux, client, conf)

	http.ListenAndServe(":8080", wrappedMux)
}
