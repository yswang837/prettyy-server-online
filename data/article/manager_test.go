package article

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewManager(t *testing.T) {
	manager, err := NewManager()
	fmt.Println("err......", err)
	assert.Equal(t, nil, err)
	fmt.Printf("client.....%+v\n", manager.client)
	manager.Add(&Article{
		Aid:        "123123",
		Title:      "123123",
		Content:    "123123",
		CoverImg:   "",
		Summary:    "",
		ReadNum:    0,
		CommentNum: 0,
		CollectNum: 0,
		Uid:        123123,
		CreateTime: time.Time{},
		UpdateTime: time.Time{},
	})
}
