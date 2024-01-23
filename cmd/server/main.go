package main

import (
	"net/http"

	web "github.com/GoExpertCurso/GoRateLimiterFC/internal/infra/web"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", web.Home)

	http.ListenAndServe(":8080", mux)
}
