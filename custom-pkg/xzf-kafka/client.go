package xzf_kafka

import (
	"github.com/Shopify/sarama"
)

type KafkaClient struct {
	C sarama.Client
}

func NewKafkaClient(brokers []string, conf *sarama.Config) (*KafkaClient, error) {
	c, err := sarama.NewClient(brokers, conf)
	return &KafkaClient{c}, err
}

func (k *KafkaClient) Topics() ([]string, error) {
	return k.C.Topics()
}

func (k *KafkaClient) Partitions(topic string) ([]int32, error) {
	return k.C.Partitions(topic)
}

func (k *KafkaClient) GetOffset(topic string, partition int32, time int64) (int64, error) {
	return k.C.GetOffset(topic, partition, time)
}

func (k *KafkaClient) Close() error {
	return k.C.Close()
}
