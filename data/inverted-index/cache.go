package inverted_index

import (
	"prettyy-server-online/custom-pkg/xzf-redis"
)

type ManagerRedis struct {
	redis *xzf_redis.Client
}

func NewManagerRedis() (*ManagerRedis, error) {
	conn, err := xzf_redis.NewClient("inverted-index")
	if err != nil {
		return nil, err
	}
	cache := &ManagerRedis{
		redis: conn,
	}
	return cache, nil
}

func (ms *ManagerRedis) HMSet(key string, fieldsValues map[string]interface{}) (string, error) {
	return ms.redis.HMSet(key, fieldsValues)
}

func (ms *ManagerRedis) HGetAll(key string) (map[string]string, error) {
	return ms.redis.HGetAll(key)
}

func (ms *ManagerRedis) Close() error {
	return nil
}