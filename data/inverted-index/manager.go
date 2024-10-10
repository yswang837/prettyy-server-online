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
	if i == nil || i.Idx == "" || i.AttrValue == "" || i.Typ == "" {
		return tool.ErrParams
	}
	now := time.Now()
	i.CreateTime = now
	i.UpdateTime = now
	if err := m.master(buildKey(i.Typ, i.AttrValue)).Create(i).Error; err != nil {
		return err
	}
	return nil
}

func (m *Manager) Delete(typ, attrValue string, index string) error {
	if typ == "" || attrValue == "" || index == "" {
		return tool.ErrParams
	}
	return m.master(buildKey(typ, attrValue)).Scopes(withTyp(typ), withAttrValue(attrValue), withIndex(index)).Delete(&InvertedIndex{}).Error
}

func (m *Manager) Update(typ, attrValue string, index string) error {
	if typ == "" || attrValue == "" || index == "" {
		return tool.ErrParams
	}
	updateVal := map[string]interface{}{"index": index, "update_time": time.Now}
	return m.master(buildKey(typ, attrValue)).Scopes(withTyp(typ), withAttrValue(attrValue)).Updates(updateVal).Error
}

func (m *Manager) Get(typ, attrValue, index string) (i []*InvertedIndex, err error) {
	if typ == "" || attrValue == "" {
		return nil, tool.ErrParams
	}
	i = []*InvertedIndex{}
	db := m.slave(buildKey(typ, attrValue))
	if index != "" {
		db.Scopes(withIndex(index))
	}
	if err = db.Scopes(withTyp(typ), withAttrValue(attrValue)).Find(&i).Error; err != nil {
		return nil, err
	}
	if len(i) == 0 {
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

func withTyp(typ string) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Where("typ = ?", typ)
	}
}

func withIndex(index string) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Where("idx = ?", index)
	}
}
