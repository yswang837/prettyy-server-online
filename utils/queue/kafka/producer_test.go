package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
)

func TestKafkaProducer_Push(t *testing.T) {
	k, err := NewProducer("first-topic")
	require.Equal(t, nil, err)
	err = k.Start()
	require.Equal(t, nil, err)
	k.OnSuccess(func(msg *sarama.ProducerMessage) {
		m, _ := msg.Value.Encode()
		fmt.Printf("i am success. msg: %s\n", string(m))

	})
	k.OnError(func(err error) {
		fmt.Printf("sorry, push failed. err: %s\n", err.Error())

	})
	defer func() {
		_ = k.Close()
	}()
	var wg sync.WaitGroup
	succ := 0
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				err = k.Push(fmt.Sprintf("new message i: %d, j: %d", i, j))
				if err != nil {
					fmt.Printf("an error equal: %s\n", err.Error())
				} else {
					succ++
				}
			}
		}(i)
	}
	wg.Wait()
	fmt.Println("timeout cnt: ", k.timeoutCnt, "succ: ", succ)
}
