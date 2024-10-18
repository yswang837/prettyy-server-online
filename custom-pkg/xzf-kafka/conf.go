package xzf_kafka

import (
	"errors"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/pelletier/go-toml"
	"log"
	"os"
	"time"
)

var (
	partitioner = map[string]func(string) sarama.Partitioner{
		"HashPartitioner":       sarama.NewHashPartitioner,       // 哈希
		"RoundRobinPartitioner": sarama.NewRoundRobinPartitioner, // 轮询
		"RandomPartitioner":     sarama.NewRandomPartitioner,     // 随机
		"ManualPartitioner":     sarama.NewManualPartitioner,     // 手动指定固定的分区
	}
	offsetInitial = map[string]int64{
		"newest": sarama.OffsetNewest,
		"oldest": sarama.OffsetOldest,
	}

	// 将分区分配给消费者的策略
	balanceStrategy = map[string]sarama.BalanceStrategy{
		"range":      sarama.BalanceStrategyRange,      // 按照范围平均分配，如1~3分区分给消费者A，4~6分区分给消费者B
		"roundrobin": sarama.BalanceStrategyRoundRobin, // 按轮询交替方式分配，如1，3，5分配给消费者A，2，4，6分配给消费者B
	}
	rootDir = os.Getenv("SERVICE_ROOT") + "/conf/kafka"
)

type Config struct {
	*sarama.Config
	Brokers []string
	Topics  []string
	GroupId string
	Debug   bool
}

func NewConfigFromToml(cfg *toml.Tree) *Config {
	conf := &Config{
		Config: sarama.NewConfig(),
	}
	if cfg == nil {
		return conf
	}
	if v, ok := cfg.Get("brokers").([]interface{}); ok {
		conf.Brokers = make([]string, 0)
		for _, vv := range v {
			conf.Brokers = append(conf.Brokers, vv.(string))
		}
	}
	if v, ok := cfg.Get("topics").([]interface{}); ok {
		conf.Topics = make([]string, 0)
		for _, vv := range v {
			conf.Topics = append(conf.Topics, vv.(string))
		}
	}
	if v, ok := cfg.Get("group_id").(string); ok {
		conf.GroupId = v
	}
	if v, ok := cfg.Get("producer_return_successes").(bool); ok {
		conf.Producer.Return.Successes = v
	}
	if v, ok := cfg.Get("producer_return_errors").(bool); ok {
		conf.Producer.Return.Errors = v
	}
	if v, ok := cfg.Get("producer_max_message_bytes").(int64); ok {
		conf.Producer.MaxMessageBytes = int(v)
	}
	if v, ok := cfg.Get("producer_required_acks").(int64); ok {
		conf.Producer.RequiredAcks = sarama.RequiredAcks(v)
	}
	if v, ok := cfg.Get("producer_timeout").(int64); ok {
		conf.Producer.Timeout = time.Duration(v) * time.Second
	}
	if v, ok := cfg.Get("producer_partitioner").(string); ok {
		if p, ok := partitioner[v]; ok {
			conf.Producer.Partitioner = p
		}
	}
	if v, ok := cfg.Get("producer_retry_max").(int64); ok {
		conf.Producer.Retry.Max = int(v)
	}
	if v, ok := cfg.Get("producer_retry_backoff").(int64); ok {
		conf.Producer.Retry.Backoff = time.Duration(v) * time.Millisecond
	}

	if v, ok := cfg.Get("consumer_fetch_min").(int64); ok {
		conf.Consumer.Fetch.Min = int32(v)
	}
	if v, ok := cfg.Get("consumer_fetch_max").(int64); ok {
		conf.Consumer.Fetch.Max = int32(v)
	}
	if v, ok := cfg.Get("consumer_fetch_default").(int64); ok {
		conf.Consumer.Fetch.Default = int32(v)
	}
	if v, ok := cfg.Get("consumer_return_errors").(bool); ok {
		conf.Consumer.Return.Errors = v
	}
	if v, ok := cfg.Get("consumer_offsets_inital").(string); ok {
		if i, ok := offsetInitial[v]; ok {
			conf.Consumer.Offsets.Initial = i
		}
	}
	if v, ok := cfg.Get("consumer_offsets_retry_max").(int64); ok {
		conf.Consumer.Offsets.Retry.Max = int(v)
	}

	if v, ok := cfg.Get("consumer_group_rebalance_strategy").(string); ok {
		if b, ok := balanceStrategy[v]; ok {
			conf.Consumer.Group.Rebalance.Strategy = b
		}
	}
	if v, ok := cfg.Get("version").(string); ok {
		if version, err := sarama.ParseKafkaVersion(v); err == nil {
			conf.Config.Version = version
		}
	}
	if v, ok := cfg.Get("client_id").(string); ok {
		conf.Config.ClientID = v
	}
	if v, ok := cfg.Get("net_sasl_enable").(bool); ok {
		conf.Net.SASL.Enable = v
	}
	if conf.Net.SASL.Enable {
		if v, ok := cfg.Get("net_sasl_mechanism").(string); ok {
			conf.Net.SASL.Mechanism = sarama.SASLMechanism(v)
		}
		if v, ok := cfg.Get("net_sasl_user").(string); ok {
			conf.Net.SASL.User = v
		}
		if v, ok := cfg.Get("net_sasl_password").(string); ok {
			conf.Net.SASL.Password = v
		}
	}

	if v, ok := cfg.Get("debug").(bool); ok {
		conf.Debug = v
	}
	return conf
}

// NewConfigFromFile 根据配置文件的路径解析出配置
// 参数 file: 配置文件的绝对/相对路径
// 该配置文件的格式需要为 toml 格式
func NewConfigFromFile(file string) *Config {
	tree, err := toml.LoadFile(file)
	if err != nil {
		log.Panic(err)
	}
	return NewConfigFromToml(tree)
}

// NewProducerConfigByNameGroup 根据默认的配置路径解析出配置
// 参数 name: 目录名，目录下包含该kafka在不同group下的不同配置，默认使用default.conf
// 参数 group: 配置文件名。
// 配置文件的路径解析方式为：os.Getenv("SERVICE_ROOT") + "/conf/kafka/producer/" + name + "/" + group + ".toml"
func NewProducerConfigByNameGroup(name string, group string) *Config {
	return newConfigByNameGroup(name, group, joinProducerPath)
}

// NewConsumerConfigByNameGroup 根据默认的配置路径解析出配置
// 参数 name: 目录名，目录下包含该kafka在不同group下的不同配置，默认使用default.conf
// 参数 group: 配置文件名。
// 配置文件的路径解析方式为：os.Getenv("SERVICE_ROOT") + "/conf/kafka/consumer/" + name + "/" + group + ".toml"
func NewConsumerConfigByNameGroup(name string, group string) *Config {
	return newConfigByNameGroup(name, group, joinConsumerPath)
}

func newConfigByNameGroup(name string, group string, joinPath func(name, group string) string) *Config {
	cfgFile := joinPath(name, group)
	filename, err := checkFile(cfgFile)
	if err != nil {
		cfgFile = joinPath(name, "default")
		filename, err = checkFile(cfgFile)
	}
	if err != nil {
		log.Panic(err)
	}
	return NewConfigFromFile(filename)
}

func joinConsumerPath(name string, group string) string {
	return fmt.Sprintf("%s/consumer/%s/%s.toml", rootDir, name, group)
}

func joinProducerPath(name string, group string) string {
	return fmt.Sprintf("%s/producer/%s/%s.toml", rootDir, name, group)
}
func checkFile(filename string) (string, error) {
	fi, err := os.Stat(filename)
	if err != nil && !os.IsExist(err) {
		return "", err
	}
	if fi.IsDir() || fi.Size() == 0 {
		return "", errors.New("config: invalid contents")
	}
	return filename, nil
}
