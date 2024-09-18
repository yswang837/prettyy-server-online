package xzf_redis

import (
	"github.com/garyburd/redigo/redis"
	"gopkg.in/ini.v1"
)

type Client struct {
	Config *ini.File
	Pool   *redis.Pool
}

func NewClient(filename string) (*Client, error) {
	cfg, err := NewConfigByName(filename)
	if err != nil {
		return nil, err
	}
	c := &Client{}
	c.Pool = &redis.Pool{
		MaxIdle:     cfg.MaxIdle,
		MaxActive:   cfg.MaxActive,
		IdleTimeout: cfg.IdleTimeout,
		Wait:        cfg.Wait,
		Dial: func() (redis.Conn, error) {
			client, err := redis.Dial("tcp",
				cfg.Master[0],
				redis.DialConnectTimeout(cfg.ConnTimeout),
				redis.DialReadTimeout(cfg.ReadTimeout),
				redis.DialWriteTimeout(cfg.WriteTimeout))
			if err != nil {
				return nil, err
			}
			if _, err = client.Do("AUTH", cfg.Password); err != nil {
				return nil, err
			}
			return client, nil
		},
	}
	_, err = c.Pool.Get().Do("ping")
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Client) Setex(key string, expiration int, value interface{}) (string, error) {
	return redis.String(c.Pool.Get().Do("SETEX", key, expiration, value))
}

func (c *Client) Get(key string) (string, error) {
	return redis.String(c.Pool.Get().Do("GET", key))
}

func (c *Client) HMSet(key string, fieldsValues map[string]interface{}) (string, error) {
	args := []interface{}{key}
	// 遍历 map，将字段和值添加到参数切片中
	for k, v := range fieldsValues {
		args = append(args, k, v)
	}
	return redis.String(c.Pool.Get().Do("HMSET", args...))
}

func (c *Client) HSet(key string, field string, value string) (int64, error) {
	return redis.Int64(c.Pool.Get().Do("HSET", key, field, value))
}

func (c *Client) Del(key string) (uint64, error) {
	return redis.Uint64(c.Pool.Get().Do("DEL", key))
}

func (c *Client) HGetAll(key string) (map[string]string, error) {
	return redis.StringMap(c.Pool.Get().Do("HGETALL", key))
}

func (c *Client) Incr(key string) (int64, error) {
	return redis.Int64(c.Pool.Get().Do("INCR", key))
}

func (c *Client) Set(key string, value string) (string, error) {
	return redis.String(c.Pool.Get().Do("SET", key, value))
}
