package main

import (
	"fmt"
	"net/http"

	"github.com/GoExpertCurso/GoRateLimiterFC/configs"
	data "github.com/GoExpertCurso/GoRateLimiterFC/internal/infra/db"
	web "github.com/GoExpertCurso/GoRateLimiterFC/internal/infra/web"
	mid "github.com/GoExpertCurso/GoRateLimiterFC/internal/infra/web/middleware"
)

func main() {
	config, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	conf := configs.NewConf(config.TokenLimit, config.IPLimit)

	redisStrategy := &data.RedisDatabaseStrategy{}
	dbClient := data.NewDatabaseClient(redisStrategy)
	redisClient, err := dbClient.Connect(config)
	if err != nil {
		fmt.Println(err)
		return
	}

	dbClient.SetClient(redisClient)
	defer dbClient.Disconnect()

	mux := http.NewServeMux()
	mux.HandleFunc("/", web.Home)

	wrappedMux := mid.RateLimitMiddleware(mux, dbClient, conf)

	http.ListenAndServe(":8080", wrappedMux)
}
