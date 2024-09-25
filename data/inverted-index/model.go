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

// typ值说明
// 1: email -> uid
// 2: uid -> aid

// InvertedIndex 面向数据库，联合主键(attr_value, number)
type InvertedIndex struct {
	Typ        string    `json:"typ"`         // 见typ值说明
	AttrValue  string    `json:"attr_value"`  // 见typ值说明
	Index      string    `json:"index"`       // 索引值
	CreateTime time.Time `json:"create_time"` // 创建时间
	UpdateTime time.Time `json:"update_time"` // 更新时间，可用于换绑邮箱，换绑手机号等，预留功能
}

func BuildPrimaryKey(Typ, AttrValue string) string {
	if Typ == "" || AttrValue == "" {
		return ""
	}
	return Typ + AttrValue
}

func (i *InvertedIndex) TableName() string {
	return tablePrefix + tool.Crc(BuildPrimaryKey(i.Typ, i.AttrValue), tableNum)
}

func (i *InvertedIndex) String() string {
	b, _ := json.Marshal(i)
	return string(b)
}
