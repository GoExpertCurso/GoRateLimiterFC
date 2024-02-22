package ratelimiter

import (
	"fmt"
	"log"
	"strconv"
	"time"

	entity "github.com/GoExpertCurso/GoRateLimiterFC/internal/entity"
	"github.com/go-redis/redis"
)

type RateLimiter struct {
	client *redis.Client
	limit  int
	window time.Duration
}

func NewRateLimiter(client *redis.Client, limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		client: client,
		limit:  limit,
		window: window,
	}
}

func (rl *RateLimiter) Allow(ipOrToken string, request entity.Request) bool {
	key := ipOrToken
	now := time.Now().Unix()

	res, _ := rl.client.TxPipelined(func(p redis.Pipeliner) error {
		p.HIncrBy(key, "count", 1)
		p.Expire(key, rl.window)
		p.HGet(key, "timestamp")
		return nil
	})

	if len(res) < 3 {
		fmt.Println("Resposta do pipeline menor do que o esperado")
		return false
	}

	count := rl.getCount(ipOrToken)

	timestamp, _ := res[2].(*redis.StringCmd).Int64()

	if now-int64(rl.window.Seconds()) > timestamp {
		rl.client.HSet(key, "timestamp", now)
		rl.client.HSet(key, "count", 1)
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

	fmt.Println()
	fmt.Println("Count:", count)
	fmt.Println("Limit:", limit)
	fmt.Println("Block:", count > limit)
	log.Println()
	fmt.Println()

	return count > limit
}

func (rl *RateLimiter) getCount(ipOrToken string) int64 {
	hashKey := ipOrToken
	field := "count"

	// Contexto para operações Redis
	//ctx := context.Background()

	// Obter o valor do campo no hash
	val, err := rl.client.HGet(hashKey, field).Result()
	if err == redis.Nil {
		fmt.Println("Campo não encontrado")
	} else if err != nil {
		fmt.Println("Erro:", err)
	}

	result, _ := strconv.ParseInt(val, 10, 64)

	return result
}
