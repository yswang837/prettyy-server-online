package xzf_redis

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSet(t *testing.T) {
	c, err := NewClient("user")
	assert.Equal(t, nil, err)
	b, err := c.Setex("1714113169@qq.com", 300, 392922)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, b)
	code, err := c.Get("1714113169@qq.com")
	assert.Equal(t, nil, err)
	fmt.Println(code)
	b, err = c.HMSet("aaa", map[string]interface{}{"bb": "cc"})
	fmt.Println(err)
	fmt.Println(b)
}
