package ratelimiter

import (
	"fmt"
	"strconv"
	"time"

	"github.com/GoExpertCurso/GoRateLimiterFC/internal/infra/db"
	"github.com/go-redis/redis"
)

type RateLimiter struct {
	client db.DatabaseClient
	limit  int
	window time.Duration
}

func NewRateLimiter(client *db.DatabaseClient, limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		client: *client,
		limit:  limit,
		window: window,
	}
}

func (rl *RateLimiter) Allow(ipOrToken string) bool {
	key := ipOrToken
	now := time.Now().Unix()
	res, _ := rl.client.PipelineTX(key, rl.window)

	if len(res.([]redis.Cmder)) < 3 {
		fmt.Println("Resposta do pipeline menor do que o esperado")
		return false
	}

	count := rl.getCount(key)

	timestamp, _ := res.([]redis.Cmder)[2].(*redis.StringCmd).Int64()

	if now-int64(rl.window.Seconds()) > timestamp {
		rl.client.Set(key, "timestamp", now)
		rl.client.Set(key, "count", 1)
		return true
	}

	if count > int64(rl.limit) {
		return false
	}

	return true
}

func (rl *RateLimiter) Block(ipOrToken string, limit int64) bool {
	key := ipOrToken
	count := rl.getCount(key)
	return count > limit
}

func (rl *RateLimiter) getCount(ipOrToken string) int64 {
	hashKey := ipOrToken
	field := "count"

	val, err := rl.client.Get(hashKey, field)
	if err == redis.Nil {
		fmt.Println("Campo não encontrado")
	} else if err != nil {
		fmt.Println("Erro:", err)
	}

	result, _ := strconv.ParseInt(val.(string), 10, 64)

	return result
}
