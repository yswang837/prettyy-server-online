package article

import (
	"encoding/json"
	"prettyy-server-online/utils/tool"
	"time"
)

const (
	tableNum    = 2
	tablePrefix = "article_"
)

// Article 面向数据库
type Article struct {
	Aid        string    `json:"aid"`         // 文章id，雪花算法生成
	Title      string    `json:"title"`       // 文章标题
	Content    string    `json:"content"`     // 文章内容
	CoverImg   string    `json:"cover_img"`   // 文章封面url
	Summary    string    `json:"summary"`     // 文章摘要
	ReadNum    int       `json:"read_num"`    // 文章的阅读数
	CommentNum int       `json:"comment_num"` // 文章的评论数
	CollectNum int       `json:"collect_num"` // 文章的收藏数
	Uid        int64     `json:"uid"`         // 本篇文章属于哪个作者
	CreateTime time.Time `json:"create_time"` // 文章的创建时间
	UpdateTime time.Time `json:"update_time"` // 文章的更新时间
}

func (a *Article) TableName() string {
	return tablePrefix + tool.Crc(a.Aid, tableNum)
}

func (a *Article) String() string {
	b, _ := json.Marshal(a)
	return string(b)
}
