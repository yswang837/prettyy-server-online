package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/pelletier/go-toml"
	"os"
	xzfKafka "prettyy-server-online/custom-pkg/xzf-kafka"
	kafkaProducer "prettyy-server-online/custom-pkg/xzf-kafka/producer"
	"prettyy-server-online/utils/metrics"
	"sync/atomic"
	"time"
)

var rootDir = fmt.Sprintf("%s/kafka/", os.Getenv("PRETTYY_CONF_ROOT"))

var (
	ErrDataType = fmt.Errorf("unsupported type")
	ErrTimeout  = fmt.Errorf("timeout")
)

const (
	defaultBufferLength = 1000
	defaultTimeout      = 10
	defaultWorkerNum    = 1
)

// Producer 基于Kafka的生产者
type Producer struct {
	producer   *kafkaProducer.Producer
	ch         chan string
	workerNum  int
	bufferLen  int
	name       string // kafka名字
	timeout    int64  // 超时时间。单位：毫秒
	timeoutCnt int64
}

// NewProducer ...
func NewProducer(name string) (*Producer, error) {
	tree, err := toml.LoadFile(fmt.Sprintf("%s%s.conf", rootDir, name))
	if err != nil {
		return nil, err
	}

	result := &Producer{
		bufferLen: defaultBufferLength,
		workerNum: defaultWorkerNum,
		timeout:   defaultTimeout,
		name:      name,
	}
	moduleCfg := tree.Get("kafka").(*toml.Tree)
	result.producer = kafkaProducer.NewProducer(xzfKafka.NewConfigFromToml(moduleCfg))

	if v, ok := tree.Get("buffer_length").(int64); ok {
		result.bufferLen = int(v)
	}

	if v, ok := tree.Get("worker_num").(int64); ok {
		result.workerNum = int(v)
	}

	if v, ok := tree.Get("timeout").(int64); ok {
		result.timeout = v
	}

	result.OnError(result.getDefaultErrorHandler())
	return result, nil
}

func (k *Producer) SetWorkerNum(workerNum int) *Producer {
	k.workerNum = workerNum
	return k
}

func (k *Producer) SetBufferLength(bufferLength int) *Producer {
	k.bufferLen = bufferLength
	return k
}

func (k *Producer) OnSuccess(handler func(*sarama.ProducerMessage)) {
	k.producer.OnSuccess(handler)
}

func (k *Producer) OnError(handler func(err error)) {
	k.producer.OnError(handler)
}

func (k *Producer) Start() error {
	k.ch = make(chan string, k.bufferLen)
	k.producer.Produce(k.ch)
	return nil
}

// Push ...
func (k *Producer) Push(message interface{}) (err error) {
	defer func() {
		v := metrics.Success
		switch err {
		case nil:
			v = metrics.Success
		case ErrDataType:
			v = metrics.ErrDatatype
		case ErrTimeout:
			v = fmt.Sprintf("%d_%s", k.timeout, metrics.Timeout)
		}
		metrics.KafkaCounter.Inc(&metrics.KafkaCounterTags{Name: k.name, Type: v})
	}()
	msg, ok := message.(string)
	if !ok {
		err = ErrDataType
		return
	}

	select {
	case k.ch <- msg:
		return
	default:
	}

	timer := time.NewTimer(time.Duration(k.timeout) * time.Millisecond)
	defer timer.Stop()
	select {
	case k.ch <- msg:
		return
	case <-timer.C:
		atomic.AddInt64(&k.timeoutCnt, 1)
		err = ErrTimeout
		return
	}
}

// Close ...
func (k *Producer) Close() error {
	if err := k.producer.Close(); err != nil {
		return err
	}
	close(k.ch)
	return nil
}

func (k *Producer) getDefaultErrorHandler() func(err error) {
	return func(err error) {
		metrics.KafkaCounter.Inc(&metrics.KafkaCounterTags{Name: k.name, Type: metrics.ErrSend})
	}
}
