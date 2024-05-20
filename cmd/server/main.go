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
	configs.LoadEnv()
	redisStrategy := &data.RedisDatabaseStrategy{}
	dbClient := data.NewDatabaseClient(redisStrategy)
	redisClient, err := dbClient.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}

	dbClient.SetClient(redisClient)
	defer dbClient.Disconnect()

	mux := http.NewServeMux()
	mux.HandleFunc("/", web.Home)

	wrappedMux := mid.RateLimitMiddleware(mux, dbClient)

	http.ListenAndServe(":8080", wrappedMux)
}
