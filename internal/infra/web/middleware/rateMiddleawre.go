package web

import (
	"fmt"
	"net"
	"net/http"
	"time"

	entity "github.com/GoExpertCurso/GoRateLimiterFC/internal/entity"
	rl "github.com/GoExpertCurso/GoRateLimiterFC/pkg/rateLimiter"
	"github.com/go-redis/redis"
)

func RateLimitMiddleware(next http.Handler, client *redis.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := getToken(r)
		ip, err := getIp()
		if err != nil {
			panic(err)
		}

		request := entity.NewRequest(*ip, token)
		request.LimitCheck()

		limiter := rl.NewRateLimiter(client, int(request.Limit), time.Second)
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
	fmt.Printf("API: %s\n", api)
	fmt.Printf("TOKEN: %s\n", token)
	if token != "" {
		return token
	}
	return api
}
