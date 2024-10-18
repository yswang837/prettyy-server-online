package producer

import (
	"github.com/Shopify/sarama"
	"log"
	xzfKafka "prettyy-server-online/custom-pkg/xzf-kafka"
	"sync"
)

type producer interface {
	init(*xzfKafka.Config) error
	send(*sarama.ProducerMessage)
	close() error
	handleSuccess(func(*sarama.ProducerMessage))
	handleError(func(err error))
}

type Producer struct {
	producer       // 将接口嵌入到结构体中，从而直接拥有了接口的方法，init,send等
	conf           *xzfKafka.Config
	wg             *sync.WaitGroup
	wg2            *sync.WaitGroup
	successHandler func(*sarama.ProducerMessage) // 函数类型
	errorHandler   func(error)                   // 函数类型
}

func NewProducer(conf *xzfKafka.Config) *Producer {
	var producer producer
	if conf.Config.Producer.RequiredAcks == sarama.NoResponse {
		producer = &asyncProducer{}
	} else {
		producer = &syncProducer{err: make(chan error, 100)}
	}
	return &Producer{
		producer: producer,
		conf:     conf,
		wg:       &sync.WaitGroup{},
		wg2:      &sync.WaitGroup{},

		successHandler: defaultSuccessHandler,
		errorHandler:   defaultErrorHandler,
	}
}

func (p *Producer) OnSuccess(handler func(*sarama.ProducerMessage)) {
	p.successHandler = handler
}

func (p *Producer) OnError(handler func(err error)) {
	p.errorHandler = handler
}

func (p *Producer) Produce(messages chan string) {
	if err := p.producer.init(p.conf); err != nil {
		log.Panic(err)
	}
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		for msg := range messages {
			pmsg := &sarama.ProducerMessage{
				Topic: p.conf.Topics[0],
				Key:   nil,
				Value: sarama.StringEncoder(msg),
			}
			p.producer.send(pmsg)
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
