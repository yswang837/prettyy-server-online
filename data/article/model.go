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
	Tags       string    `json:"tags"`        // 文章标签，以逗号分隔，最多10个标签，由用户发文的时候打标签
	Visibility string    `json:"visibility"`  // 文章的可见性, 1-全部可见 2-VIP可见 3-粉丝可见 4-仅我可见
	Typ        string    `json:"typ"`         // 文章类型，1-原创 2-转载 3-翻译
	ShareNum   int       `json:"share_num"`   // 文章的分享数
	CommentNum int       `json:"comment_num"` // 文章的评论数
	LikeNum    int       `json:"like_num"`    // 文章的点赞数
	ReadNum    int       `json:"read_num"`    // 文章的阅读数
	CollectNum int       `json:"collect_num"` // 文章的收藏数
	Status     string    `json:"status"`      // 文章的状态，1-正常 2-审核中 3-审核不通过/草稿箱 4-已删除
	Uid        int64     `json:"uid"`         // 文章属于哪个作者
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
