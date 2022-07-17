package cachedmysql

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v9"
)

type cache struct {
	rds    *redis.Client
	expiry time.Duration
}

func (c cache) Del(keys ...string) error {
	return c.DelCtx(context.Background(), keys...)
}

func (c cache) DelCtx(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}

	if _, err := c.rds.Del(ctx, keys...).Result(); err != nil {
		return err
	}

	return nil
}

func (c cache) Get(key string, val any, query func(val interface{}) error) error {
	return c.GetWithExpireCtx(context.Background(), key, val, func(v any, expire time.Duration) error {
		return query(val)
	}, c.expiry)
}

func (c cache) GetCtx(ctx context.Context, key string, val any, query func(val interface{}) error) error {
	return c.GetWithExpireCtx(ctx, key, val, func(v any, expire time.Duration) error {
		return query(val)
	}, c.expiry)
}

func (c cache) GetWithExpire(key string, val any,
	query func(v any, expire time.Duration) error, expire time.Duration) error {
	return c.GetWithExpireCtx(context.Background(), key, val, query, expire)
}

func (c cache) GetWithExpireCtx(ctx context.Context, key string, val any,
	query func(v any, expire time.Duration) error, expire time.Duration) error {

	data, err := sf.do(key, func() (any, error) {
		data, err := c.rds.Get(ctx, key).Result()
		if err != nil || data == "" {
			if err := query(val, expire); err != nil {
				return val, err
			}
		}
		err = json.Unmarshal([]byte(data), val)
		return val, err
	})
	if err == nil {
		val = data
	}
	return err
}

func (c cache) Set(key string, val any) error {
	return c.SetWithExpireCtx(context.Background(), key, val, c.expiry)
}

func (c cache) SetCtx(ctx context.Context, key string, val any) error {
	return c.SetWithExpireCtx(ctx, key, val, c.expiry)
}

func (c cache) SetWithExpire(key string, val any, expire time.Duration) error {
	return c.SetWithExpireCtx(context.Background(), key, val, expire)
}

func (c cache) SetWithExpireCtx(ctx context.Context, key string, val any, expire time.Duration) error {
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}

	_, err = c.rds.Set(ctx, key, string(data), expire).Result()
	return err
}
