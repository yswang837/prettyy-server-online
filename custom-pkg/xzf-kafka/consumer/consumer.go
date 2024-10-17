package consumer

import (
	"github.com/Shopify/sarama"
	"log"
	xzfKafka "prettyy-server-online/custom-pkg/xzf-kafka"
	"sync"
)

type Actuator interface {
	Init() error
	Do(message *sarama.ConsumerMessage) error
	Destroy() error
}

type consumer interface {
	init(*xzfKafka.Config, Actuator) error
	consume()
	close() error
}

type Handler struct {
	actuator Actuator
	wg       sync.WaitGroup
	debug    bool
}

func (h *Handler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Handler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Handler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/master/consumer_group.go#L27-L29s
	for message := range claim.Messages() {
		//log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s, offset = %d, partition = %d", string(message.Value), message.Timestamp, message.Topic, message.Offset, message.Partition)
		if err := h.actuator.Do(message); err != nil {
			log.Printf("Message claimed err: topic = %s, partition = %d, offset = %d, key = %s, value = %s, err = %s", message.Topic, message.Partition, message.Offset, message.Key, string(message.Value), err)
		}
		if !h.debug {
			session.MarkMessage(message, "")
		}
	}
	return nil
}

type Consumer struct {
	conf     *xzfKafka.Config
	actuator Actuator
	consumer consumer
}

func NewConsumer(conf *xzfKafka.Config, act Actuator) *Consumer {
	var consumer consumer
	//sarama ConsumerGroup 要求kafka版本号不低于V0_10_2_0
	if conf.Version.IsAtLeast(sarama.V0_10_2_0) {
		consumer = &highConsumerGroup{}
	} else {
		consumer = &simpleConsumerGroup{}
	}
	if err := consumer.init(conf, act); err != nil {
		log.Panic(err)
	}
	return &Consumer{
		conf:     conf,
		actuator: act,
		consumer: consumer,
	}
}

func (c *Consumer) Consume() {
	c.consumer.consume()
}

func (c *Consumer) Close() error {
	return c.consumer.close()
}
