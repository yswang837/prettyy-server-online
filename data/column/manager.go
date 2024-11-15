package column

import (
	"errors"
	"gorm.io/gorm"
	"os"
	xzfMysql "prettyy-server-online/custom-pkg/xzf-mysql"
	"prettyy-server-online/utils/tool"
	"time"
)

type Manager struct {
	client *xzfMysql.Client
}

func NewManager() (*Manager, error) {
	cfg, err := xzfMysql.NewConfig("column")
	if err != nil {
		return nil, err
	}
	client, err := xzfMysql.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return &Manager{client: client}, nil
}

func (m *Manager) Add(c *Column) error {
	if c == nil || c.Uid == 0 || c.Cid == "" || c.Title == "" {
		return tool.ErrParams
	}
	now := time.Now()
	c.CreateTime = now
	c.UpdateTime = now
	return m.master(c.Cid).Create(c).Error
}

func (m *Manager) Get(cid string) (*Column, error) {
	if cid == "" {
		return nil, tool.ErrParams
	}
	c := &Column{}
	if err := m.slave(cid).Scopes(withCid(cid)).First(c).Error; err != nil {
		return nil, err
	}
	if c.Cid == "" {
		return nil, errors.New("record not found")
	}
	return c, nil
}

func (m *Manager) master(key string) *gorm.DB {
	return m.client.Master().Model(&Column{}).Scopes(selectTable(key))
}

func (m *Manager) slave(key string) *gorm.DB {
	return m.client.Slave().Model(&Column{}).Scopes(selectTable(key))
}

func selectTable(key string) func(tx *gorm.DB) *gorm.DB {
	if os.Getenv("idc") == "dev" {
		return func(tx *gorm.DB) *gorm.DB {
			return tx.Table(tablePrefix + "0")
		}
	}
	num := tool.Crc(key, tableNum)
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Table(tablePrefix + num)
	}
}

func withCid(cid string) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Where("cid = ?", cid)
	}
}
