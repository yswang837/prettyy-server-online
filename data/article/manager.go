package article

import (
	"errors"
	"gorm.io/gorm"
	"math/rand"
	"os"
	xzfMysql "prettyy-server-online/custom-pkg/xzf-mysql"
	"prettyy-server-online/utils/tool"
	"strconv"
	"time"
)

type Manager struct {
	client *xzfMysql.Client
}

func NewManager() (*Manager, error) {
	cfg, err := xzfMysql.NewConfig("article")
	if err != nil {
		return nil, err
	}
	client, err := xzfMysql.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return &Manager{client: client}, nil
}

func (m *Manager) Add(a *Article) error {
	if a == nil || a.Aid == "" || a.Uid == 0 || a.Content == "" || a.Title == "" {
		return tool.ErrParams
	}
	now := time.Now()
	a.CreateTime, a.UpdateTime = now, now
	if err := m.master(a.Aid).Create(a).Error; err != nil {
		return err
	}
	return nil
}

func (m *Manager) Get(aid string) (*Article, error) {
	if aid == "" {
		return nil, tool.ErrParams
	}
	a := &Article{}
	if err := m.slave(aid).Scopes(withAid(aid)).First(a).Error; err != nil {
		return nil, err
	}
	if a.Aid == "" {
		return nil, errors.New("record not found")
	}
	return a, nil
}

// GetArticleList 简单查询则参数传递对应类型的零值，也支持分页查询
func (m *Manager) GetArticleList(uid int64, page, pageSize int) (art []*Article, err error) {
	art = []*Article{}
	db := m.slave(strconv.Itoa(rand.Intn(100))) // 随机从从库中找一个表获取数据，它不是aid
	if uid >= 10000 {
		db.Scopes(withUid(strconv.FormatInt(uid, 10)))
	}
	if pageSize > 0 {
		db.Limit(pageSize)
	} else {
		db.Limit(50)
	}
	offset := 0
	if page > 1 && pageSize > 0 {
		offset = (page - 1) * pageSize
	}
	if err = db.Offset(offset).Find(&art).Error; err != nil {
		return nil, err
	}

	if len(art) == 0 {
		return nil, errors.New("record not found")
	}
	return
}

func (m *Manager) Delete(aid string) error {
	if aid == "" {
		return tool.ErrParams
	}
	a := &Article{}
	if err := m.master(aid).Scopes(withAid(aid)).Delete(a).Error; err != nil {
		return err
	}
	return nil
}

func (m *Manager) master(aid string) *gorm.DB {
	return m.client.Master().Model(&Article{}).Scopes(selectTable(aid))
}

func (m *Manager) slave(aid string) *gorm.DB {
	return m.client.Slave().Model(&Article{}).Scopes(selectTable(aid))
}

func selectTable(aid string) func(tx *gorm.DB) *gorm.DB {
	if os.Getenv("PRETTYY_TEST") == "dev" {
		return func(tx *gorm.DB) *gorm.DB {
			return tx.Table(tablePrefix + "0")
		}
	}
	num := tool.Crc(aid, tableNum)
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Table(tablePrefix + num)
	}
}

func withAid(aid string) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Where("aid = ?", aid)
	}
}

func withUid(uid string) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Where("uid = ?", uid)
	}
}
