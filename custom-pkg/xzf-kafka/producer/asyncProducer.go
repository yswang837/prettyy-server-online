package producer

import (
	"github.com/Shopify/sarama"
	xzfKafka "prettyy-server-online/custom-pkg/xzf-kafka"
)

type asyncProducer struct {
	sarama.AsyncProducer
}

func (p *asyncProducer) init(conf *xzfKafka.Config) error {
	producer, err := sarama.NewAsyncProducer(conf.Brokers, conf.Config)
	if err != nil {
		return err
	}
	p.AsyncProducer = producer
	return nil
}

func (p *asyncProducer) send(msg *sarama.ProducerMessage) {
	p.AsyncProducer.Input() <- msg
}

func (p *asyncProducer) close() error {
	return p.Close() // 底层逻辑保证了会把积压的消息消费完
}

func (p *asyncProducer) handleSuccess(f func(msg *sarama.ProducerMessage)) {
	for msg := range p.Successes() {
		f(msg)
	}
}

func (p *asyncProducer) handleError(f func(err error)) {
	for err := range p.Errors() {
		f(err)
	}
}
