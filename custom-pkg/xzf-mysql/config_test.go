package xzf_mysql

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewConfig(t *testing.T) {
	cfg, err := NewConfig("user")
	assert.Equal(t, nil, err)
	fmt.Printf("cfg\n%+v\n", cfg)

}
