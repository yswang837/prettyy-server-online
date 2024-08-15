package xzf_redis

import (
	"errors"
	"fmt"
	"gopkg.in/ini.v1"
	"os"
	"strings"
	"time"
)

var (
	rootDir            = "config/redis"
	ErrInvalidContents = errors.New("config: invalid contents")
	defaultGroup       = "default"
	redisSection       = "redis"
)

// Config ...
type Config struct {
	MaxIdle      int
	MaxActive    int
	Wait         bool
	IdleTimeout  time.Duration
	ConnTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	Master       []string
	Slave        []string
	Name         string
}

func DefaultConfig() *Config {
	return &Config{
		MaxIdle:      64,
		MaxActive:    64,
		Wait:         true,
		IdleTimeout:  30 * time.Second,
		ConnTimeout:  100 * time.Millisecond,
		ReadTimeout:  100 * time.Millisecond,
		WriteTimeout: 100 * time.Millisecond,
		Master:       []string{},
		Slave:        []string{},
	}
}

func NewConfigByName(name string) (*Config, error) {
	filename, err := checkName(name)
	if err != nil {
		return nil, err
	}
	cfg, err := NewConfigByFile(filename)
	if err != nil {
		return nil, err
	}
	cfg.Name = name
	return cfg, nil
}

func NewConfigByFile(configFile string) (*Config, error) {
	if configFile == "" {
		return nil, errors.New("redis config name is empty")
	}
	f, err := ini.Load(configFile)
	if err != nil {
		return nil, err
	}
	cfg := DefaultConfig()
	if maxIdle, _ := f.Section(redisSection).Key("max_idle").Int(); maxIdle != 0 {
		cfg.MaxIdle = maxIdle
	}
	if maxActive, _ := f.Section(redisSection).Key("max_active").Int(); maxActive != 0 {
		cfg.MaxActive = maxActive
	}
	if wait, _ := f.Section(redisSection).Key("wait").Bool(); wait {
		cfg.Wait = wait
	}
	if idleTimeout, _ := f.Section(redisSection).Key("idle_timeout").Int(); idleTimeout != 0 {
		cfg.IdleTimeout = time.Duration(idleTimeout)
	}
	if readTimeout, _ := f.Section(redisSection).Key("read_timeout").Int(); readTimeout != 0 {
		cfg.ReadTimeout = time.Duration(readTimeout)
	}
	if writeTimeout, _ := f.Section(redisSection).Key("write_timeout").Int(); writeTimeout != 0 {
		cfg.WriteTimeout = time.Duration(writeTimeout)
	}
	if connTimeout, _ := f.Section(redisSection).Key("conn_timeout").Int(); connTimeout != 0 {
		cfg.ConnTimeout = time.Duration(connTimeout)
	}
	cfg.Master = strings.Split(f.Section(redisSection).Key("master").String(), ",")
	cfg.Slave = strings.Split(f.Section(redisSection).Key("slave").String(), ",")
	return cfg, nil
}

func checkName(name string) (filename string, err error) {
	cfgFile := fmt.Sprintf("%s/%s/%s.ini", rootDir, name, defaultGroup)
	filename, err = checkFile(cfgFile)

	return
}

func checkFile(filename string) (string, error) {
	fi, err := os.Stat(filename)
	if err != nil && !os.IsExist(err) {
		return "", err
	}
	if fi.IsDir() || fi.Size() == 0 {
		return "", ErrInvalidContents
	}
	return filename, nil
}
