package xzf_kafka

import (
	"github.com/IBM/sarama"
	"github.com/pelletier/go-toml"
	"time"
)

var (
	partitioner = map[string]func(string) sarama.Partitioner{
		"HashPartitioner":       sarama.NewHashPartitioner,
		"RoundRobinPartitioner": sarama.NewRoundRobinPartitioner,
		"RandomPartitioner":     sarama.NewRandomPartitioner,
		"ManualPartitioner":     sarama.NewManualPartitioner,
	}
	offsetInitial = map[string]int64{
		"newest": sarama.OffsetNewest,
		"oldest": sarama.OffsetOldest,
	}
	balanceStrategy = map[string]sarama.BalanceStrategy{
		"range":      sarama.BalanceStrategyRange,
		"roundrobin": sarama.BalanceStrategyRoundRobin,
	}
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
