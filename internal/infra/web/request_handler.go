package web

import (
	"log"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s - %d", r.Method, http.StatusOK)
	w.Write([]byte("Initializing Rate Limiter Solution...."))
}
