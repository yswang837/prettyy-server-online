package queue

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestNewProducer(t *testing.T) {
	err := os.Setenv("PRODUCER_FILE_DIR", "/Users/xzf")
	require.Equal(t, err, nil)
	producer, err := NewProducer("test")
	require.Equal(t, err, nil)
	defer func() {
		_ = producer.Close()
	}()
	for i := 0; i < 2e4; i++ {
		_ = producer.Push(fmt.Sprintf("message: %d", i))
	}
}
