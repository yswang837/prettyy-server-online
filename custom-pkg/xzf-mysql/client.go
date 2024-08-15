package xzf_mysql

import (
	"database/sql"
	"errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"strings"
	"time"
)

type Client struct {
	config *Config
	master []*gorm.DB // 连接池
	slave  []*gorm.DB // 连接池
}

func NewClient(cfg *Config) (*Client, error) {
	if cfg == nil {
		return nil, errors.New("mysql new client err: params is nil")
	}
	c := &Client{config: cfg}
	dsn, err := cfg.makeDsn()
	if err != nil {
		return nil, err
	}
	if c.master, err = c.newPool(dsn.master); err != nil {
		return nil, err
	}
	if c.slave, err = c.newPool(dsn.slave); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Client) newPool(dsn []string) (db []*gorm.DB, err error) {
	if len(dsn) == 0 {
		return nil, errors.New("newPool params err")
	}
	db = make([]*gorm.DB, len(dsn))
	for i, addr := range dsn {
		if db[i], err = c.newConn(addr); err != nil {
			return nil, err
		}
	}
	return db, err
}

func (c *Client) newConn(dsn string) (*gorm.DB, error) {
	d := mysql.Open(dsn)
	opt := &gorm.Config{
		Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			Colorful:                  strings.ToLower(os.Getenv("MYSQL_LOG_COLOR")) == "on",
			IgnoreRecordNotFoundError: true,
			LogLevel:                  12,
		}),
	}
	if c.config.Logger != nil {
		opt.Logger = c.config.Logger
	}
	if c.config.LogLevel > 0 {
		opt.Logger = opt.Logger.LogMode(logger.LogLevel(c.config.LogLevel))
	}
	db, err := gorm.Open(d, opt)
	if err != nil {
		return nil, err
	}
	pool, err := db.DB()
	if err != nil {
		return nil, err
	}
	pool.SetConnMaxIdleTime(time.Duration(c.config.MaxIdleConn))
	pool.SetMaxOpenConns(c.config.MaxOpenConn)
	pool.SetConnMaxLifetime(c.config.MaxLifeTime)
	return db, nil
}

func (c *Client) Master() *gorm.DB {
	if len(c.master) > 0 {
		return c.master[0]
	}
	return nil
}

func (c *Client) Slave() *gorm.DB {
	if len(c.slave) > 0 {
		return c.slave[0]
	}
	return nil
}

func (c *Client) SetMaster(username, password string, hosts []string) error {
	if username == "" || password == "" || len(hosts) <= 0 {
		return errors.New("params err")
	}
	c.config.Master.Username = username
	c.config.Master.Password = password
	c.config.Master.Hosts = hosts
	return nil
}

func (c *Client) SetSlave(username, password string, hosts []string) error {
	if username == "" || password == "" || len(hosts) <= 0 {
		return errors.New("params err")
	}
	c.config.Slave.Username = username
	c.config.Slave.Password = password
	c.config.Slave.Hosts = hosts
	return nil
}

func (c *Client) closePool(pool []*gorm.DB) (err error) {
	var db *sql.DB
	for k, p := range pool {
		if db, err = p.DB(); err != nil {
			break
		}
		if err = db.Close(); err != nil {
			break
		}
		pool[k] = nil
	}
	return
}

func (c *Client) Close() error {
	if err := c.closePool(c.master); err != nil {
		return err
	}
	if err := c.closePool(c.slave); err != nil {
		return err
	}
	c.master = nil
	c.slave = nil
	return nil
}
