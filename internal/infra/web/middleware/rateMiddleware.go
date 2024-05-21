package web

import (
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	entity "github.com/GoExpertCurso/GoRateLimiterFC/internal/entity"
	"github.com/GoExpertCurso/GoRateLimiterFC/internal/infra/db"
	rl "github.com/GoExpertCurso/GoRateLimiterFC/pkg/rateLimiter"
)

func RateLimitMiddleware(next http.Handler, client interface{}) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := getToken(r)
		ip, err := getIp()
		if err != nil {
			panic(err)
		}

		key := getKey(*ip, token)
		limit, duration := getLimitAndDuration(token)

		request := entity.NewRequest(key, limit)

		limiter := rl.NewRateLimiter(client.(*db.DatabaseClient), int(request.Limit), duration)
		allowed := limiter.Allow(key)

		if limiter.Block(*ip, int64(request.Limit)) || !allowed {
			log.Printf("Request not allowed\n")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("you have reached the maximum number of requests or actions allowed within a certain time frame"))
			log.Printf("%s - %d - you have reached the maximum number of requests or actions allowed within a certain time frame", r.Method, http.StatusTooManyRequests)
			return
		} else {
			log.Printf("Request allowed\n")
			w.WriteHeader(http.StatusOK)
			log.Printf("%s - %d - %s", r.Method, http.StatusOK, key)
			next.ServeHTTP(w, r)
		}

	})
}

func getIp() (*string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Printf("Erro ao obter o endereço IP: %v", err)
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

func getKey(ip, token string) string {
	if token != "" {
		return token
	}
	return ip
}

func getLimitAndDuration(token string) (int, time.Duration) {
	if token != "" {
		tokenLimit, _ := strconv.Atoi(os.Getenv("TOKEN_RATE_LIMIT"))
		tokenBlockDuration, _ := strconv.Atoi(os.Getenv("TOKEN_BLOCK_DURATION"))
		return tokenLimit, time.Duration(tokenBlockDuration) * time.Second
	}
	ipLimit, _ := strconv.Atoi(os.Getenv("IP_RATE_LIMIT"))
	ipBlockDuration, _ := strconv.Atoi(os.Getenv("IP_BLOCK_DURATION"))
	return ipLimit, time.Duration(ipBlockDuration) * time.Second
}
