package main

import "net/http"

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Rate Limiter Solution"))
	})

	http.ListenAndServe(":8080", mux)
}
