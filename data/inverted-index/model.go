package inverted_index

import (
	"encoding/json"
	"prettyy-server-online/utils/tool"
	"time"
)

const (
	tableNum              = 2
	tablePrefix           = "inverted_index_"
	TypEmailUid           = "1" // 通过email查uid
	TypUidAid             = "2" // 通过uid查aid
	TypUidCid             = "3" // 通过uid查cid
	TypMuidLikeSuidAid    = "4" // 通过muid,aid查suid 点赞类型
	TypMuidCollectSuidAid = "5" // 通过muid,aid查suid 收藏类型
)

// InvertedIndex 面向数据库
type InvertedIndex struct {
	Typ        string    `json:"typ"`         // 类型，将上方常量定义
	AttrValue  string    `json:"attr_value"`  // 属性值，将上方常量定义
	Idx        string    `json:"idx"`         // 想要查询的值，这里不能用index，会和sql关键字冲突
	CreateTime time.Time `json:"create_time"` // 创建时间
	UpdateTime time.Time `json:"update_time"` // 更新时间，可用于换绑邮箱，换绑手机号等，预留功能
}

// buildKey 通过typ和attr_value作为索引表名的依据
func buildKey(Typ, AttrValue string) string {
	if Typ == "" || AttrValue == "" {
		return ""
	}
	return Typ + AttrValue
}

func (i *InvertedIndex) TableName() string {
	return tablePrefix + tool.Crc(buildKey(i.Typ, i.AttrValue), tableNum)
}

func (i *InvertedIndex) String() string {
	b, _ := json.Marshal(i)
	return string(b)
}
