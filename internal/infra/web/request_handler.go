package web

import (
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Initializing Rate Limiter Solution...."))
}
