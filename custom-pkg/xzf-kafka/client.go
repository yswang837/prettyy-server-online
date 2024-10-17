package xzf_kafka

import (
	"github.com/IBM/sarama"
)

type KafkaClient struct {
	c sarama.Client
}

func NewKafkaClient(brokers []string, config *sarama.Config) (*KafkaClient, error) {
	c, err := sarama.NewClient(brokers, config)
	return &KafkaClient{c: c}, err
}

func (k *KafkaClient) Topics() ([]string, error) {
	return k.c.Topics()
}

func (k *KafkaClient) Partitions(topic string) ([]int32, error) {
	return k.c.Partitions(topic)
}

func (k *KafkaClient) GetOffset(topic string, partition int32, time int64) (int64, error) {
	return k.c.GetOffset(topic, partition, time)
}

func (k *KafkaClient) Close() error {
	return k.c.Close()
}
