package web

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/GoExpertCurso/GoRateLimiterFC/configs"
	entity "github.com/GoExpertCurso/GoRateLimiterFC/internal/entity"
	"github.com/GoExpertCurso/GoRateLimiterFC/internal/infra/db"
	rl "github.com/GoExpertCurso/GoRateLimiterFC/pkg/rateLimiter"
)

func RateLimitMiddleware(next http.Handler, client interface{}, limits *configs.Conf) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := getToken(r)
		ip, err := getIp()
		if err != nil {
			panic(err)
		}

		request := entity.NewRequest(*ip, token, *limits)
		request.LimitCheck()

		limiter := rl.NewRateLimiter(client.(*db.DatabaseClient), int(request.Limit), time.Second)
		allowed := limiter.Allow(chooseType(*ip, token), *request)

		if limiter.Block(*ip, int64(request.Limit)) || !allowed {
			fmt.Printf("Request not allowed\n")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("Too Many Requests"))
			return
		} else {
			fmt.Printf("Request allowed\n")
			w.WriteHeader(http.StatusOK)
			next.ServeHTTP(w, r)
		}

	})
}

func getIp() (*string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		fmt.Println("Erro ao obter o endereço IP:", err)
		return nil, err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	ip := localAddr.IP.String()
	return &ip, nil
}

func getToken(request *http.Request) string {
	return request.Header.Get("api-key")
}

func chooseType(api string, token string) string {
	if token != "" {
		return token
	}
	return api
}
