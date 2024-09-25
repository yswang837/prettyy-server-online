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

// GetHomeArticleList 主页文章列表数据，简单查询则参数传递对应类型的零值，也支持分页查询，也支持条件查询
func (m *Manager) GetHomeArticleList(page, pageSize int, visibility string) (art []*Article, err error) {
	art = []*Article{}
	db := m.slave(strconv.Itoa(rand.Intn(100))) // 随机从从库中找一个表获取数据，它不是aid
	if visibility != "" {
		db.Scopes(withVisibility(visibility))
	}
	if pageSize > 0 {
		db.Limit(pageSize)
	} else {
		pageSize = 20
		db.Limit(20)
	}
	offset := 0
	if page > 1 {
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

// GetContentManageArticleList 内容管理页面文章列表数据，简单查询则参数传递对应类型的零值，也支持分页查询，也支持条件查询
func (m *Manager) GetContentManageArticleList(aids []string, visibility, typ string) (artList []*Article, count int64, err error) {
	artList = []*Article{}
	db := m.client.Slave().Model(&Article{})
	if visibility != "" {
		db.Scopes(withVisibility(visibility))
	}
	if typ != "" {
		db.Scopes(withTyp(typ))
	}
	for _, aid := range aids {
		db = m.slave(aid)
		art := &Article{}
		if err = db.Scopes(withAid(aid)).First(art).Error; err != nil {
			continue
		}
		count++
	}
	if len(artList) == 0 {
		return nil, 0, errors.New("record not found")
	}
	return
}

func (m *Manager) Delete(aid string, uid int64) error {
	if aid == "" {
		return tool.ErrParams
	}
	a := &Article{Aid: aid, Uid: uid}
	if err := m.master(aid).Scopes(withAid(aid), withUid(strconv.FormatInt(uid, 10))).Delete(a).Error; err != nil {
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
	if os.Getenv("idc") == "dev" {
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

func withVisibility(visibility string) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Where("visibility = ?", visibility)
	}
}

func withTyp(tpy string) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Where("typ = ?", tpy)
	}
}

func withAids(aids []string) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Where("aid in ?", aids)
	}
}
