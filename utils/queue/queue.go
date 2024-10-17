package queue

import (
	"fmt"
	"os"
	"prettyy-server-online/utils/queue/file"
	"prettyy-server-online/utils/queue/kafka"
	"strconv"
)

// Consumer ...
type Consumer interface {
	Consume(func(interface{}) error) error
	Close() error
}

// Producer ...
type Producer interface {
	Push(interface{}) error
	Close() error
}

// NewProducer ... 如果配置了PRODUCER_FILE_DIR环境变量，则会将数据写到该环境变量下的文件中；否则写到kafka里
func NewProducer(name string) (Producer, error) {
	dir := os.Getenv("PRODUCER_FILE_DIR")
	if dir != "" {
		fileProducer, err := file.NewProducer(fmt.Sprintf("%s/%s.txt", dir, name))
		if err != nil {
			return nil, err
		}
		buf, num := os.Getenv("PRODUCER_FILE_BUFFER"), os.Getenv("PRODUCER_FILE_WORKERNUM")
		if buffer, _ := strconv.Atoi(buf); buffer > 0 {
			fileProducer.SetBufferLength(buffer)
		}
		if workerNum, _ := strconv.Atoi(num); workerNum > 0 {
			fileProducer.SetWorkerNum(workerNum)
		}
		_ = fileProducer.Start()
		return fileProducer, nil
	}
	kafkaProducer, err := kafka.NewProducer(name)
	if err != nil {
		return nil, err
	}
	_ = kafkaProducer.Start()
	return kafkaProducer, nil
}
