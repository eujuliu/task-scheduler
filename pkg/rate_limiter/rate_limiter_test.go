package ratelimiter_test

import (
	"context"
	ratelimiter "scheduler/pkg/rate_limiter"
	"scheduler/pkg/redis"
	"testing"

	"github.com/go-redis/redismock/v9"

	. "scheduler/test"
)

func TestSlidingWindowCounterLimiter(t *testing.T) {
	rdb := &redis.Redis{}
	db, mock := redismock.NewClientMock()
	ctx := context.Background()

	err := SetPrivateField(rdb, "db", db)
	Ok(t, err)

	limiter := ratelimiter.NewSlidingWindowCounterLimiter(rdb, 5, 60, 20)

	key := "rate_limit:client_1"

	mock.ExpectHGetAll(key).SetVal(map[string]string{"test": "5"})
	_, allowed := limiter.Allowed(ctx, "client_1")
	Equals(t, false, allowed)

	Ok(t, mock.ExpectationsWereMet())
}
