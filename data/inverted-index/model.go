package user

import (
	"encoding/json"
	"prettyy-server-online/utils/tool"
	"time"
)

const (
	tableNum    = 2
	tablePrefix = "inverted_index_"
)

// InvertedIndex 面向数据库，联合主键(attr_value, number)
type InvertedIndex struct {
	AttrValue  string    `json:"attr_value"`  // 属性值，目前是：email的值，后续可以新增phone的值
	Number     string    `json:"number"`      // 当number为1时，表示email，当number为2时，表示phone
	Uid        int64     `json:"uid"`         // 用户id
	CreateTime time.Time `json:"create_time"` // 创建时间
}

func BuildPrimaryKey(AttrValue, Number string) string {
	if AttrValue == "" || Number == "" {
		return ""
	}
	return AttrValue + Number
}

func (i *InvertedIndex) TableName() string {
	return tablePrefix + tool.Crc(BuildPrimaryKey(i.AttrValue, i.Number), tableNum)
}

func (i *InvertedIndex) String() string {
	b, _ := json.Marshal(i)
	return string(b)
}
