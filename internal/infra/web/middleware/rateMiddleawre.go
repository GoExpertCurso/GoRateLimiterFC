package web

import (
	"fmt"
	"log"
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

		key, _, timeD := chooseType(*ip, token)

		duration := time.Duration(int64(timeD) * int64(time.Second))

		limiter := rl.NewRateLimiter(client.(*db.DatabaseClient), int(request.Limit), duration)
		allowed := limiter.Allow(key, *request)

		if limiter.Block(*ip, int64(request.Limit)) || !allowed {
			fmt.Printf("Request not allowed\n")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("you have reached the maximum number of requests or actions allowed within a certain time frame"))
			log.Printf("%s - %d - you have reached the maximum number of requests or actions allowed within a certain time frame", r.Method, http.StatusTooManyRequests)
			return
		} else {
			fmt.Printf("Request allowed\n")
			w.WriteHeader(http.StatusOK)
			log.Printf("%s - %d", r.Method, http.StatusOK)
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

func chooseType(api string, token string) (string, int, int) {
	if token != "" {
		limit, time := configs.GetTokenLimit()
		return token, limit, time
	}
	limit, time := configs.GetIpLimit()
	return api, limit, time
}
