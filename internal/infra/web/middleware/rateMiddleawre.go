package web

import (
	"fmt"
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
		fmt.Println("Key:", key)
		fmt.Println("Duration: ", duration)
		fmt.Println("Limit: ", limit)

		request := entity.NewRequest(key, limit)
		//request.LimitCheck()

		//duration = time.Duration(duration * time.Second)

		limiter := rl.NewRateLimiter(client.(*db.DatabaseClient), int(request.Limit), duration)
		allowed := limiter.Allow(key)

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
