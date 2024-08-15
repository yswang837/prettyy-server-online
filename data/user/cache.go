package user

import (
	"prettyy-server-online/custom-pkg/xzf-redis"
)

type ManagerRedis struct {
	redis *xzf_redis.Client
}

func NewManagerRedis() (*ManagerRedis, error) {
	conn, err := xzf_redis.NewClient("user")
	if err != nil {
		return nil, err
	}
	cache := &ManagerRedis{
		redis: conn,
	}
	return cache, nil
}

func (ms *ManagerRedis) SetWithTTL(key string, expire int, val string) (string, error) {
	return ms.redis.Setex(key, expire, val)
}

func (ms *ManagerRedis) Incr(key string) (int64, error) {
	uid, _ := ms.redis.Get(key)
	if uid == "" {
		// uid从10000开始自增
		_, _ = ms.redis.Set(key, "10000")
	}
	return ms.redis.Incr(key)
}

func (ms *ManagerRedis) HMSet(key string, fieldsValues map[string]interface{}) (string, error) {
	return ms.redis.HMSet(key, fieldsValues)
}

func (ms *ManagerRedis) HSet(key, field, value string) (int64, error) {
	return ms.redis.HSet(key, field, value)
}

func (ms *ManagerRedis) HGetAll(key string) (map[string]string, error) {
	return ms.redis.HGetAll(key)
}

func (ms *ManagerRedis) Get(key string) (string, error) {
	return ms.redis.Get(key)
}

func (ms *ManagerRedis) Close() error {
	return nil
}
