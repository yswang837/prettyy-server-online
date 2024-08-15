package xzf_mysql

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
)

func TestNewClient(t *testing.T) {
	cfg, err := NewConfig("user")
	assert.Equal(t, nil, err)
	fmt.Println(cfg)
	c, err := NewClient(cfg)
	assert.Equal(t, nil, err)
	fmt.Println(c)
	m := func() *gorm.DB {
		return c.Master().Exec("use aaa")
	}
	sql := "create table a1(a int,b int)"
	m().Exec(sql)
}
