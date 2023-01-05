package store

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/wgo-admin/backend/internal/pkg/errno"
)

type ICacheStore interface {
	SetStruct(ctx context.Context, key string, obj interface{}, seconds int) error
	GetStruct(ctx context.Context, key string) (obj interface{}, err error)
	Del(ctx context.Context, key ...string) error
}

var _ ICacheStore = (*cache)(nil)

func newCache(rdb *redis.Client) *cache {
	return &cache{rdb}
}

type cache struct {
	rdb *redis.Client
}

func (c *cache) Del(ctx context.Context, key ...string) error {
	return c.rdb.Del(ctx, key...).Err()
}

// GetStruct 获取一个 struct
func (c *cache) GetStruct(ctx context.Context, key string) (obj interface{}, err error) {
	result, err := c.rdb.Do(ctx, "GET", key).Result()

	err = json.Unmarshal([]byte(fmt.Sprintf("%v", result)), &obj)
	if err != nil {
		return nil, errno.ErrSerialization
	}

	return obj, nil
}

// SetStruct 缓存一个 struct
func (c *cache) SetStruct(ctx context.Context, key string, obj interface{}, seconds int) error {
	// 序列化成 json
	json, err := json.Marshal(obj)
	if err != nil {
		return errno.ErrSerialization
	}

	// seconds 等于0则不设置过期时间
	if seconds == 0 {
		err = c.rdb.Do(ctx, "SET", key, string(json)).Err()
	} else {
		err = c.rdb.Do(ctx, "SET", key, string(json), "EX", seconds).Err()
	}

	if err != nil {
		return err
	}

	return nil
}
