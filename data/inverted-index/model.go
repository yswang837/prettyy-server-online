package inverted_index

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
	Typ        string    `json:"typ"`         // 当number为1时，表示email，当number为2时，表示phone
	AttrValue  string    `json:"attr_value"`  // 属性值，目前是：email的值，后续可以新增phone的值
	Index      string    `json:"index"`       // 索引值
	CreateTime time.Time `json:"create_time"` // 创建时间
	UpdateTime time.Time `json:"update_time"` // 更新时间，可用于换绑邮箱，换绑手机号等，预留功能
}

func BuildPrimaryKey(AttrValue, Number string) string {
	if AttrValue == "" || Number == "" {
		return ""
	}
	return AttrValue + Number
}

func (i *InvertedIndex) TableName() string {
	return tablePrefix + tool.Crc(BuildPrimaryKey(i.AttrValue, i.Typ), tableNum)
}

func (i *InvertedIndex) String() string {
	b, _ := json.Marshal(i)
	return string(b)
}
