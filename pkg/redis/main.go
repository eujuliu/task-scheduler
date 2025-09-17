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
	ctx    context.Context
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
		ctx:    context.Background(),
	}
}

func (r *Redis) BeginTransaction() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.tx = r.db.TxPipeline()
}

func (r *Redis) ExecTransaction() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.tx == nil {
		return fmt.Errorf("you need to initialize the transaction first")
	}

	_, err := r.tx.Exec(r.ctx)
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

func (r *Redis) HGetAll(key string) (map[string]string, error) {
	return r.db.HGetAll(r.ctx, key).Result()
}

func (r *Redis) HIncrBy(key, field string, incr int64) (int64, error) {
	var tx redis.Cmdable = r.db

	if r.tx != nil {
		tx = r.tx
	}

	return tx.HIncrBy(r.ctx, key, field, incr).Result()
}

func (r *Redis) HExpire(
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

	return tx.HExpireWithArgs(r.ctx, key, expiration, args, fields...).Result()
}
