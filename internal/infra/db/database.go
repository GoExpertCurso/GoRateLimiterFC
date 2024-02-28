package db

import (
	"time"

	"github.com/GoExpertCurso/GoRateLimiterFC/configs"
)

type DatabaseStrategy interface {
	Connect(*configs.Conf) (interface{}, error)
	Disconnect() error
	Get(hashkey, field string) (any, error)
	Set(string, string, int64)
	PipelineTX(key string, time time.Duration) (interface{}, error)
	Delete(key string) error
	SetClient(client interface{}) error
}

type DatabaseClient struct {
	strategy DatabaseStrategy
}

func NewDatabaseClient(strategy DatabaseStrategy) *DatabaseClient {
	return &DatabaseClient{strategy: strategy}
}

func (c *DatabaseClient) UseDatabaseStrategy(strategy DatabaseStrategy) {
	c.strategy = strategy
}

func (c *DatabaseClient) Connect(config *configs.Conf) (interface{}, error) {
	return c.strategy.Connect(config)
}

func (c *DatabaseClient) Disconnect() error {
	return c.strategy.Disconnect()
}

func (c *DatabaseClient) Get(hashKey, field string) (interface{}, error) {
	return c.strategy.Get(hashKey, field)
}

func (c *DatabaseClient) Set(key string, value string, param int64) {
	c.strategy.Set(key, value, param)
}

func (c *DatabaseClient) PipelineTX(key string, time time.Duration) (interface{}, error) {
	return c.strategy.PipelineTX(key, time)
}

func (c *DatabaseClient) Delete(key string) error {
	return c.strategy.Delete(key)
}

func (c *DatabaseClient) SetClient(client interface{}) error {
	return c.strategy.SetClient(client)
}
