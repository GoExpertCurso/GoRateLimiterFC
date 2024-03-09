package db

import (
	"errors"
	"log"
	"time"

	"github.com/GoExpertCurso/GoRateLimiterFC/configs"
	"github.com/go-redis/redis"
)

type RedisDatabaseStrategy struct {
	client      *redis.Client
	isConnected bool
}

func (r *RedisDatabaseStrategy) Connect(config *configs.Conf) (any, error) {
	r.isConnected = true
	client := redis.NewClient(&redis.Options{
		Addr:     config.DBHost + ":" + config.DBPort,
		Password: config.DBPassword,
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	log.Printf("Connected to Redis")
	return client, nil
}

func (r *RedisDatabaseStrategy) Disconnect() error {
	if r.isConnected {
		r.client.Close()
		r.isConnected = false
		log.Println("Desconectado do Redis")
	}
	return nil
}

func (r *RedisDatabaseStrategy) Get(hashKey, field string) (any, error) {
	return r.client.HGet(hashKey, field).Result()
}

func (r *RedisDatabaseStrategy) Set(key string, value string, param int64) {
	r.client.HSet(key, value, param)
}

func (r *RedisDatabaseStrategy) PipelineTX(key string, window time.Duration) (any, error) {
	res, err := r.client.TxPipelined(func(p redis.Pipeliner) error {
		p.HIncrBy(key, "count", 1)
		p.Expire(key, window)
		p.HGet(key, "timestamp")
		return nil
	})
	return res, err
}

func (r *RedisDatabaseStrategy) Delete(key string) error {
	err := r.client.Del(key).Err()
	if err != nil {
		log.Printf("Error deleting key: %s", err)
		return err
	}
	return nil
}

func (r *RedisDatabaseStrategy) SetClient(client interface{}) error {
	if client == nil {
		return errors.New("client is nil")
	}
	r.client = client.(*redis.Client)
	return nil
}
