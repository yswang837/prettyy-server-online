package user

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewManager(t *testing.T) {
	manager, err := NewManager()
	assert.Equal(t, nil, err)
	fmt.Printf("%+v\n", manager.client)
	manager.Add(&User{Uid: 11111})
}
