package column

import (
	"encoding/json"
	"prettyy-server-online/utils/tool"
	"time"
)

const (
	tableNum    = 2
	tablePrefix = "column_"
)

// Column 面向数据库
type Column struct {
	Cid          string    `json:"cid"`            // 专栏id，雪花算法生成
	Title        string    `json:"title"`          // 专栏标题
	CoverImg     string    `json:"cover_img"`      // 专栏封面url
	Summary      string    `json:"summary"`        // 专栏摘要
	FrontDisplay string    `json:"front_display"`  // 是否前台展示，1-展示 2-不展示
	IsFreeColumn string    `json:"is_free_column"` // 是否免费专栏，1-免费 2-付费
	SubscribeNum int       `json:"subscribe_num"`  // 专栏订阅数
	Uid          int64     `json:"uid"`            // 专栏属于哪个作者
	CreateTime   time.Time `json:"create_time"`    // 专栏的创建时间
	UpdateTime   time.Time `json:"update_time"`    // 专栏的更新时间
}

func (c *Column) TableName() string {
	return tablePrefix + tool.Crc(c.Cid, tableNum)
}

func (c *Column) String() string {
	b, _ := json.Marshal(c)
	return string(b)
}
