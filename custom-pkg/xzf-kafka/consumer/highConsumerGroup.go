package consumer

import (
	"context"
	"github.com/Shopify/sarama"
	"log"
	xzfKafka "prettyy-server-online/custom-pkg/xzf-kafka"
)

type highConsumerGroup struct {
	sarama.ConsumerGroup
	handler sarama.ConsumerGroupHandler
	topics  []string
	closeCh chan struct{}
	ctx     context.Context
	cancel  context.CancelFunc
}

func (h *highConsumerGroup) init(conf *xzfKafka.Config, act Actuator) error {
	client, err := xzfKafka.NewKafkaClient(conf.Brokers, conf.Config)
	if err != nil {
		return err
	}
	consumer, err := sarama.NewConsumerGroupFromClient(conf.GroupId, client.C)
	if err != nil {
		return err
	}
	if err := act.Init(); err != nil {
		return err
	}
	h.handler = &Handler{actuator: act, debug: conf.Debug}
	h.ConsumerGroup = consumer
	h.topics = conf.Topics
	h.closeCh = make(chan struct{})
	h.ctx, h.cancel = context.WithCancel(context.Background())
	return nil
}

func (h *highConsumerGroup) consume() {
	var err error
loop:
	for {
		select {
		case <-h.closeCh:
			break loop
		default:
		}
		// `Consume` should be called inside an infinite loop, when a
		// server-side rebalance happens, the consumer session will need to be
		// recreated to get the new claims
		if err = h.ConsumerGroup.Consume(h.ctx, h.topics, h.handler); err != nil {
			log.Printf("Error from consumer: %v", err)
			break
		}
	}
	// consume 函数退出之后，actuator可能需要通过 AfterExit 做一些清理工作。
	if afterExit, ok := h.handler.(*Handler).actuator.(interface{ AfterExit(err error) }); ok {
		afterExit.AfterExit(err)
	}
}

func (h *highConsumerGroup) close() error {
	h.cancel()
	h.closeCh <- struct{}{}
	return h.Close()
}
