package file

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
)

// Producer 基于文件的生产者
type Producer struct {
	f         *os.File
	total     int64
	success   int64
	failed    int64
	ch        chan string
	wg        *sync.WaitGroup
	workerNum int
	bufferLen int
}

// NewProducer ...
func NewProducer(path string) (*Producer, error) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return nil, err
	}
	result := &Producer{
		f:         f,
		workerNum: 1,
		wg:        &sync.WaitGroup{},
	}
	return result, nil
}

func (f *Producer) SetWorkerNum(workerNum int) *Producer {
	f.workerNum = workerNum
	return f
}

func (f *Producer) SetBufferLength(bufferLength int) *Producer {
	f.bufferLen = bufferLength
	return f
}

func (f *Producer) Start() error {
	f.ch = make(chan string, f.bufferLen)
	for i := 0; i < f.workerNum; i++ {
		f.wg.Add(1)
		go func() {
			defer f.wg.Done()
			f.worker()
		}()
	}
	return nil
}

func (f *Producer) worker() {
	for msg := range f.ch {
		atomic.AddInt64(&f.total, 1)
		_, err := f.f.WriteString(fmt.Sprintf("%s\n", msg))
		if err != nil {
			atomic.AddInt64(&f.failed, 1)
			continue
		}
		atomic.AddInt64(&f.success, 1)
	}
}

// Push ...
func (f *Producer) Push(msg interface{}) error {
	s, ok := msg.(string)
	if !ok {
		return errors.New("unsupported type")
	}
	f.ch <- s
	return nil
}

// Close ...
func (f *Producer) Close() error {
	close(f.ch)
	f.wg.Wait()
	_, err := fmt.Fprintf(os.Stdout, "name: %s, total: %d, success: %d, failed: %d\n", f.f.Name(), f.total, f.success, f.failed)
	if err != nil {
		return err
	}
	err = f.f.Close()
	return err
}
