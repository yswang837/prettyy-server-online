package producer

import (
	"github.com/IBM/sarama"
	xzfKafka "prettyy-server-online/custom-pkg/xzf-kafka"
)

type syncProducer struct {
	sarama.SyncProducer
	err chan error
}

func (s *syncProducer) init(config *xzfKafka.Config) error {
	p, err := sarama.NewSyncProducer(config.Brokers, config.Config)
	if err != nil {
		return err
	}
	s.SyncProducer = p
	return nil
}

func (s *syncProducer) send(message *sarama.ProducerMessage) {
	defer func() {
		recover()
	}()
	_, _, err := s.SyncProducer.SendMessage(message)
	s.err <- err
}

func (s *syncProducer) handleSuccess(func(*sarama.ProducerMessage)) {}

func (s *syncProducer) handleError(f func(err error)) {
	for err := range s.err {
		f(err)
	}
}

func (s *syncProducer) close() error {
	close(s.err)
	return s.Close()
}
