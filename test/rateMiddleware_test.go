package test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/GoExpertCurso/GoRateLimiterFC/configs"
	data "github.com/GoExpertCurso/GoRateLimiterFC/internal/infra/db"
	web "github.com/GoExpertCurso/GoRateLimiterFC/internal/infra/web/middleware"
	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis"
)

func setupMockRedis(t *testing.T) (*data.DatabaseClient, *miniredis.Miniredis) {
	mredis, err := miniredis.Run()
	if err != nil {
		t.Fatalf("could not start miniredis: %v", err)
	}
	dbClient := data.NewDatabaseClient(&data.RedisDatabaseStrategy{})
	rdb := redis.NewClient(&redis.Options{
		Addr: mredis.Addr(),
	})
	dbClient.SetClient(rdb)
	return dbClient, mredis
}

func TestRateLimiterByIP(t *testing.T) {
	rdb, mredis := setupMockRedis(t)
	defer mredis.Close()

	os.Setenv("IP_RATE_LIMIT", "5")
	os.Setenv("IP_BLOCK_DURATION", "10")
	configs.LoadEnv()

	handler := http.NewServeMux()
	handler.Handle("/", web.RateLimitMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}), rdb))

	req := httptest.NewRequest("GET", "https://www.google.com.br/", nil)
	req.RemoteAddr = "192.168.1.1"

	for i := 0; i < 5; i++ {
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected status OK but got %v", w.Code)
		}
	}

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("expected status 429 but got %v", w.Code)
	}
}

func TestRateLimiterByToken(t *testing.T) {
	rdb, mredis := setupMockRedis(t)
	defer mredis.Close()

	os.Setenv("TOKEN_RATE_LIMIT", "3")
	os.Setenv("TOKEN_BLOCK_DURATION", "10")
	configs.LoadEnv()

	handler := http.NewServeMux()
	handler.Handle("/", web.RateLimitMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}), rdb))

	req := httptest.NewRequest("GET", "https://www.google.com.br/", nil)
	req.Header.Set("api-key", "abc123")

	for i := 0; i < 3; i++ {
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected status OK but got %v", w.Code)
		}
	}

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("expected status 429 but got %v", w.Code)
	}
}

func TestRateLimiterPriority(t *testing.T) {
	rdb, mredis := setupMockRedis(t)
	defer mredis.Close()

	os.Setenv("IP_RATE_LIMIT", "2")
	os.Setenv("TOKEN_RATE_LIMIT", "3")
	os.Setenv("IP_BLOCK_DURATION", "10")
	os.Setenv("TOKEN_BLOCK_DURATION", "10")
	configs.LoadEnv()

	handler := http.NewServeMux()
	handler.Handle("/", web.RateLimitMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}), rdb))

	req := httptest.NewRequest("GET", "https://www.google.com.br/", nil)
	req.RemoteAddr = "192.168.1.1"
	req.Header.Set("api-key", "abc123")

	for i := 0; i < 3; i++ {
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected status OK but got %v", w.Code)
		}
	}

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("expected status 429 but got %v", w.Code)
	}
}

func setup() {
	configs.LoadEnv()
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}
