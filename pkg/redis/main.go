package redis

import (
	"context"
	"fmt"
	"reflect"
	"scheduler/internal/config"
	"sync"
	"time"

	redis "github.com/redis/go-redis/v9"
)

type Redis struct {
	db     *redis.Client
	config *config.RedisConfig
	mu     sync.Mutex
	tx     redis.Pipeliner
}

func NewRedis(config *config.RedisConfig) *Redis {
	db := redis.NewClient(&redis.Options{
		Addr:       config.Addr,
		Password:   config.Password,
		Username:   config.Username,
		DB:         config.DB,
		ClientName: "scheduler",
	})

	return &Redis{
		db:     db,
		config: config,
	}
}

func (r *Redis) BeginTransaction() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.tx = r.db.TxPipeline()
}

func (r *Redis) ExecTransaction(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.tx == nil {
		return fmt.Errorf("you need to initialize the transaction first")
	}

	_, err := r.tx.Exec(ctx)
	r.tx = nil

	if err != nil {
		return err
	}

	return nil
}

func (r *Redis) DiscardTransaction() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.tx == nil {
		return fmt.Errorf("you need to initialize the transaction first")
	}

	r.tx.Discard()
	r.tx = nil

	return nil
}

func (r *Redis) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return r.db.HGetAll(ctx, key).Result()
}

func (r *Redis) HIncrBy(ctx context.Context, key, field string, incr int64) (int64, error) {
	var tx redis.Cmdable = r.db

	if r.tx != nil {
		tx = r.tx
	}

	return tx.HIncrBy(ctx, key, field, incr).Result()
}

func (r *Redis) HExpire(
	ctx context.Context,
	key string,
	expiration time.Duration,
	mode string,
	fields ...string,
) ([]int64, error) {
	var tx redis.Cmdable = r.db

	if r.tx != nil {
		tx = r.tx
	}

	args := redis.HExpireArgs{}

	reflect.ValueOf(&args).Elem().FieldByName(mode).SetBool(true)

	return tx.HExpireWithArgs(ctx, key, expiration, args, fields...).Result()
}

func (r *Redis) Set(
	ctx context.Context,
	key, value string,
	expiration time.Duration,
) (string, error) {
	return r.db.Set(ctx, key, value, expiration).Result()
}

func (r *Redis) Del(
	ctx context.Context,
	key string,
) (int64, error) {
	return r.db.Del(ctx, key).Result()
}
