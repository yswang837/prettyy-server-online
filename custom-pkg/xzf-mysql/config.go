package xzf_mysql

import (
	"errors"
	"fmt"
	"gopkg.in/ini.v1"
	"gorm.io/gorm/logger"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	mysqlSection = "mysql"
	rootDir      = "config/mysql"
)

type Config struct {
	Master ConnAuth
	Slave  ConnAuth

	MaxOpenConn  int           //最大可连接数
	MaxIdleConn  int           //最大空闲连接数
	MaxLifeTime  time.Duration //连接的最大生存时间
	ConnTimeout  time.Duration //连接超时时间
	ReadTimeout  time.Duration //读超时时间
	WriteTimeout time.Duration //写超时时间

	IsParseTime bool
	Charset     string
	Location    string

	Logger   logger.Interface //gorm自带的日志包
	LogLevel int              // 10:打印所有详细日志 11:打印警告和错误日志 12:只打印错误日志 13:不打印日志

}

type ConnAuth struct {
	Username string
	Password string
	Hosts    []string //单个host的格式为ip:port:dbname
}

type dsn struct {
	master []string
	slave  []string
}

func DefaultConfig() *Config {
	return &Config{
		Master:       ConnAuth{},
		Slave:        ConnAuth{},
		MaxOpenConn:  16,
		MaxIdleConn:  8,
		MaxLifeTime:  30 * time.Minute,
		ConnTimeout:  100 * time.Microsecond,
		ReadTimeout:  100 * time.Microsecond,
		WriteTimeout: 100 * time.Microsecond,
		IsParseTime:  true,
		Charset:      "latin1",
		Location:     "Asia/Shanghai",
		Logger:       nil,
		LogLevel:     12,
	}
}

func NewConfig(name string) (*Config, error) {
	if name == "" {
		return nil, errors.New("mysql config name is empty")
	}
	configFile := fmt.Sprintf("%s/mysql/%s/default.ini", os.Getenv("PRETTYY_CONF_ROOT"), name)
	f, err := ini.Load(configFile)
	if err != nil {
		return nil, err
	}
	cfg := DefaultConfig()
	if username := f.Section(mysqlSection).Key("username").String(); username != "" {
		cfg.Master.Username, cfg.Slave.Username = username, username
	}
	if password := f.Section(mysqlSection).Key("password").String(); password != "" {
		cfg.Master.Password, cfg.Slave.Password = password, password
	}
	if maxOpenConn, _ := f.Section(mysqlSection).Key("max_open_conn").Int(); maxOpenConn != 0 {
		cfg.MaxOpenConn = maxOpenConn
	}
	if maxIdleConn, _ := f.Section(mysqlSection).Key("max_idle_conn").Int(); maxIdleConn != 0 {
		cfg.MaxIdleConn = maxIdleConn
	}
	if maxLifeTime, _ := f.Section(mysqlSection).Key("max_life_time").Int(); maxLifeTime != 0 {
		cfg.MaxLifeTime = time.Duration(maxLifeTime)
	}
	if connTimeout, _ := f.Section(mysqlSection).Key("conn_timeout").Int(); connTimeout != 0 {
		cfg.ConnTimeout = time.Duration(connTimeout)
	}
	if readTimeout, _ := f.Section(mysqlSection).Key("read_timeout").Int(); readTimeout != 0 {
		cfg.ReadTimeout = time.Duration(readTimeout)
	}
	if writeTimeout, _ := f.Section(mysqlSection).Key("write_timeout").Int(); writeTimeout != 0 {
		cfg.WriteTimeout = time.Duration(writeTimeout)
	}
	if isParseTime, _ := f.Section(mysqlSection).Key("is_parse_time").Bool(); isParseTime {
		cfg.IsParseTime = isParseTime
	}
	if charset := f.Section(mysqlSection).Key("charset").String(); charset != "" {
		cfg.Charset = charset
	}
	if location := f.Section(mysqlSection).Key("location").String(); location != "" {
		cfg.Location = location
	}
	if logLevel, _ := f.Section(mysqlSection).Key("log_level").Int(); logLevel != 0 {
		cfg.LogLevel = logLevel
	}
	cfg.Master.Hosts = strings.Split(f.Section(mysqlSection).Key("master").String(), ",")
	cfg.Slave.Hosts = strings.Split(f.Section(mysqlSection).Key("slave").String(), ",")
	return cfg, nil
}

func (c *Config) makeDsn() (*dsn, error) {

	val := url.Values{}
	val.Set("timeout", strconv.Itoa(int(c.ConnTimeout))+"s")
	val.Set("readTimeout", strconv.Itoa(int(c.ReadTimeout))+"s")
	val.Set("writeTimeout", strconv.Itoa(int(c.WriteTimeout))+"s")
	val.Set("charset", c.Charset)
	val.Set("loc", c.Location)
	if c.IsParseTime {
		val.Set("parseTime", "True")
	} else {
		val.Set("parseTime", "False")
	}
	params := val.Encode()
	toDsn := func(auth ConnAuth) ([]string, error) {
		hosts := make([]string, len(auth.Hosts))
		for i, dsn := range auth.Hosts {
			arr := strings.Split(dsn, ":")
			if len(arr) < 2 {
				return nil, fmt.Errorf("mysql hosts err: %s", dsn)
			}
			host, port := arr[0], arr[1]
			dbname := ""
			if len(arr) == 3 {
				dbname = arr[2]
			}
			hosts[i] = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", auth.Username, auth.Password, host, port, dbname, params)
		}
		return hosts, nil
	}
	dsn := &dsn{}
	var err error
	if dsn.master, err = toDsn(c.Master); err != nil {
		return nil, err
	}
	if dsn.slave, err = toDsn(c.Slave); err != nil {
		return nil, err
	}
	return dsn, nil
}
