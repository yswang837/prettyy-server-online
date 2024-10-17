package producer

import (
	"github.com/IBM/sarama"
	xzfKafka "prettyy-server-online/custom-pkg/xzf-kafka"
)

type asyncProducer struct {
	sarama.AsyncProducer
}

func (a *asyncProducer) init(config *xzfKafka.Config) error {
	ap, err := sarama.NewAsyncProducer(config.Brokers, config.Config)
	if err != nil {
		return err
	}
	a.AsyncProducer = ap
	return nil
}

func (a *asyncProducer) send(message *sarama.ProducerMessage) {
	a.AsyncProducer.Input() <- message
}

func (a *asyncProducer) handleSuccess(f func(*sarama.ProducerMessage)) {
	for msg := range a.Successes() {
		f(msg)
	}
}

func (a *asyncProducer) handleError(f func(err error)) {
	for err := range a.Errors() {
		f(err)
	}
}

func (a *asyncProducer) close() error {
	return a.Close() // sarama的底层逻辑保证了会把积压的消息消费完
}
