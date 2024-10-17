package producer

import (
	"github.com/IBM/sarama"
	"log"
	"prettyy-server-online/custom-pkg/xzf-kafka"
	"sync"
)

type producer interface {
	init(*xzf_kafka.Config) error
	send(message *sarama.ProducerMessage)
	handleSuccess(func(*sarama.ProducerMessage))
	handleError(func(err error))
	close() error
}

type Producer struct {
	producer
	conf           *xzf_kafka.Config
	wg             *sync.WaitGroup
	wg2            *sync.WaitGroup
	successHandler func(*sarama.ProducerMessage)
	errorHandler   func(error)
}

func NewProducer(conf *xzf_kafka.Config) *Producer {
	var p producer
	if conf.Config.Producer.RequiredAcks == sarama.NoResponse {
		p = &asyncProducer{}
	} else {
		p = &syncProducer{}
	}
	return &Producer{
		producer:       p,
		conf:           conf,
		wg:             &sync.WaitGroup{},
		wg2:            &sync.WaitGroup{},
		successHandler: defaultSuccessHandler,
		errorHandler:   defaultErrorHandler,
	}
}

func (p *Producer) Produce(message chan string) {
	if err := p.producer.init(p.conf); err != nil {
		log.Panic(err)
	}
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		for msg := range message {
			producerMsg := &sarama.ProducerMessage{
				Topic: p.conf.Topics[0],
				Key:   nil,
				Value: sarama.StringEncoder(msg),
			}
			p.producer.send(producerMsg)
		}
	}()
	if p.conf.Producer.Return.Errors {
		p.wg2.Add(1)
		go func() {
			defer p.wg2.Done()
			p.handleError()
		}()
	}
	if p.conf.Producer.Return.Successes {
		p.wg2.Add(1)
		go func() {
			defer p.wg2.Done()
			p.handleSuccess()
		}()
	}
}

func (p *Producer) OnSuccess(handler func(*sarama.ProducerMessage)) {
	p.successHandler = handler
}

func (p *Producer) OnError(handler func(err error)) {
	p.errorHandler = handler
}

func (p *Producer) handleSuccess() {
	p.producer.handleSuccess(p.successHandler)
}

func (p *Producer) handleError() {
	p.producer.handleError(p.errorHandler)
}

func (p *Producer) Close() error {
	p.wg.Wait() // 等到积压的消息发送完
	err := p.producer.close()
	p.wg2.Wait() // 把可能的异步的成功、错误相响应都读完
	return err
}

func defaultErrorHandler(err error) {
	log.Println(err)
}

func defaultSuccessHandler(msg *sarama.ProducerMessage) {
	log.Println(msg.Offset)
}
