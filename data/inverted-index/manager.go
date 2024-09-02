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
	if i == nil || i.Uid == 0 || i.AttrValue == "" || i.Number == "" {
		return tool.ErrParams
	}
	i.CreateTime = time.Now()
	if err := m.master(BuildPrimaryKey(i.AttrValue, i.Number)).Create(i).Error; err != nil {
		return err
	}
	return nil
}

func (m *Manager) Get(attrValue, number string) (*InvertedIndex, error) {
	if attrValue == "" || number == "" {
		return nil, tool.ErrParams
	}
	i := &InvertedIndex{}
	if err := m.slave(BuildPrimaryKey(i.AttrValue, i.Number)).Scopes(withAttrValue(attrValue), withNumber(number)).First(i).Error; err != nil {
		return nil, err
	}
	if i.Uid == 0 {
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
	if os.Getenv("PRETTYY_TEST") == "dev" {
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

func withNumber(number string) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Where("number = ?", number)
	}
}
