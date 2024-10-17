package producer

import (
	"github.com/Shopify/sarama"
	xzfKafka "prettyy-server-online/custom-pkg/xzf-kafka"
)

type syncProducer struct {
	sarama.SyncProducer
	err chan error
}

func (p *syncProducer) init(conf *xzfKafka.Config) error {
	producer, err := sarama.NewSyncProducer(conf.Brokers, conf.Config)
	if err != nil {
		return err
	}
	p.SyncProducer = producer
	return nil
}

func (p *syncProducer) send(msg *sarama.ProducerMessage) {
	defer func() {
		if recover() != nil {
		}
	}()
	_, _, err := p.SyncProducer.SendMessage(msg)
	p.err <- err
}

func (p *syncProducer) close() error {
	err := p.Close()
	close(p.err)
	return err
}

func (p *syncProducer) handleSuccess(_ func(_ *sarama.ProducerMessage)) {

}

func (p *syncProducer) handleError(f func(err error)) {
	for err := range p.err {
		f(err)
	}
}
