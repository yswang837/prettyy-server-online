package inverted_index

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
	cfg, err := xzfMysql.NewConfig("inverted-index")
	if err != nil {
		return nil, err
	}
	client, err := xzfMysql.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return &Manager{client: client}, nil
}

func (m *Manager) Add(i *InvertedIndex) error {
	if i == nil || i.Index == "" || i.AttrValue == "" || i.Typ == "" {
		return tool.ErrParams
	}
	now := time.Now()
	i.CreateTime = now
	i.UpdateTime = now
	if err := m.master(BuildPrimaryKey(i.Typ, i.AttrValue)).Create(i).Error; err != nil {
		return err
	}
	return nil
}

func (m *Manager) Get(typ, attrValue string) (*InvertedIndex, error) {
	if typ == "" || attrValue == "" {
		return nil, tool.ErrParams
	}
	i := &InvertedIndex{}
	if err := m.slave(BuildPrimaryKey(i.Typ, i.AttrValue)).Scopes(withNumber(typ), withAttrValue(attrValue)).First(i).Error; err != nil {
		return nil, err
	}
	if i.Index == "" {
		return nil, errors.New("record not found")
	}
	return i, nil
}

func (m *Manager) master(key string) *gorm.DB {
	return m.client.Master().Model(&InvertedIndex{}).Scopes(selectTable(key))
}

func (m *Manager) slave(key string) *gorm.DB {
	return m.client.Slave().Model(&InvertedIndex{}).Scopes(selectTable(key))
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

func withAttrValue(attrValue string) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Where("attr_value = ?", attrValue)
	}
}

func withNumber(typ string) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Where("typ = ?", typ)
	}
}
