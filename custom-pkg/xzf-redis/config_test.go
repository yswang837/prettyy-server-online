package xzf_redis

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewConfig(t *testing.T) {
	cfg, err := NewConfigByName("identify-code")
	assert.Equal(t, nil, err)
	fmt.Printf("cfg\n%+v\n", cfg.Master)

}
