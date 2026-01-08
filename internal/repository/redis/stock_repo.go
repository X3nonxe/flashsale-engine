package redis

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
)

type RedisStockRepository struct {
	Client *redis.Client 
}

func NewRedisStockRepository(client *redis.Client) *RedisStockRepository {
	return &RedisStockRepository{Client: client}
}

const atomicDecrScript = `
		local stock = redis.call("GET", KEYS[1])
		if not stock then return 0 end
		stock = tonumber(stock)
		local qty = tonumber(ARGV[1])
		if stock >= qty then
			redis.call("DECRBY", KEYS[1], qty)
			return 1
		end
		return 0
	`

func (r *RedisStockRepository) DecrStock(ctx context.Context, key string, qty int) error {
	res, err := r.Client.Eval(ctx, atomicDecrScript, []string{key}, qty).Result()
	
	if err != nil {
		return err
	}

	if resInt, ok := res.(int64); ok {
		if resInt == 1 {
			return nil // Success
		}
	}

	return errors.New("stock update failed: insufficient stock")
}

func (r *RedisStockRepository) SetStock(ctx context.Context, key string, qty int, expiration int) error {
	return r.Client.Set(ctx, key, qty, 0).Err()
}