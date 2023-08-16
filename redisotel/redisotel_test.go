package redisotel

import (
	"github.com/go-redis/redis/v8"
	"testing"
)

func TestTracingHook(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{})
	if rdb == nil {
		t.Errorf("failed opening connection to redis")
	}

	rdb.AddHook(TracingHook{})
}
