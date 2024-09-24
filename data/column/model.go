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

// Column 面向数据库，主键(cid)
type Column struct {
	Cid          string    `json:"cid"`            // 主键，专栏ID
	Title        string    `json:"title"`          // 专栏标题
	CoverImg     string    `json:"cover_img"`      // 专栏封面
	Summary      string    `json:"summary"`        // 专栏简介
	FrontDisplay string    `json:"front_display"`  // 专栏是否前台展示，1：展示，2：不展示
	IsFreeColumn string    `json:"is_free_column"` // 是否免费专栏，1：免费，2：收费
	SubscribeNum int64     `json:"subscribe_num"`  // 订阅人数
	Uid          int64     `json:"uid"`            // 用户id
	CreateTime   time.Time `json:"create_time"`    // 创建时间
	UpdateTime   time.Time `json:"update_time"`    // 更新时间
}

func (c *Column) TableName() string {
	return tablePrefix + tool.Crc(c.Cid, tableNum)
}

func (c *Column) String() string {
	b, _ := json.Marshal(c)
	return string(b)
}
