package user

import (
	"errors"
	"gorm.io/gorm"
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
	cfg, err := xzfMysql.NewConfig("user")
	if err != nil {
		return nil, err
	}
	client, err := xzfMysql.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return &Manager{client: client}, nil
}

func (m *Manager) Add(u *User) error {
	if u == nil || u.Uid == 0 || (u.Email == "" && u.Phone == "") {
		return tool.ErrParams
	}
	now := time.Now()
	u.CreateTime, u.UpdateTime = now, now
	if err := m.master(strconv.FormatInt(u.Uid, 10)).Create(u).Error; err != nil {
		return err
	}
	return nil
}

func (m *Manager) Get(uid string) (*User, error) {
	if uid == "" {
		return nil, tool.ErrParams
	}
	u := &User{}
	if err := m.slave(uid).Scopes(withUid(uid)).First(u).Error; err != nil {
		return nil, err
	}
	if u.Uid == 0 {
		return nil, errors.New("record not found")
	}
	return u, nil
}

func (m *Manager) UpdateLoginTime(uid string) error {
	if uid == "" {
		return tool.ErrParams
	}
	if err := m.master(uid).Scopes(withUid(uid)).Update("login_time", time.Now()).Error; err != nil {
		return err
	}
	return nil
}

func (m *Manager) UpdatePassword(uid, password string) error {
	if uid == "" || password == "" {
		return tool.ErrParams
	}
	if err := m.master(uid).Scopes(withUid(uid)).Update("password", password).Error; err != nil {
		return err
	}
	return nil
}

func (m *Manager) UpdateNickName(uid, nickName string) error {
	if uid == "" || nickName == "" {
		return tool.ErrParams
	}
	if err := m.master(uid).Scopes(withUid(uid)).Update("nick_name", nickName).Error; err != nil {
		return err
	}
	return nil
}

func (m *Manager) UpdateSummary(uid, summary string) error {
	if uid == "" || summary == "" {
		return tool.ErrParams
	}
	if err := m.master(uid).Scopes(withUid(uid)).Update("summary", summary).Error; err != nil {
		return err
	}
	return nil
}

func (m *Manager) UpdateProvinceCity(uid, provinceCity string) error {
	if uid == "" || provinceCity == "" {
		return tool.ErrParams
	}
	if err := m.master(uid).Scopes(withUid(uid)).Update("province_city", provinceCity).Error; err != nil {
		return err
	}
	return nil
}

func (m *Manager) UpdateBirthdayCity(uid, birthday string) error {
	if uid == "" || birthday == "" {
		return tool.ErrParams
	}
	if err := m.master(uid).Scopes(withUid(uid)).Update("birthday", birthday).Error; err != nil {
		return err
	}
	return nil
}

func (m *Manager) UpdateGender(uid, gender string) error {
	if uid == "" || gender == "保密" || gender == "" {
		return tool.ErrParams
	}
	if err := m.master(uid).Scopes(withUid(uid)).Update("gender", gender).Error; err != nil {
		return err
	}
	return nil
}

func (m *Manager) master(uid string) *gorm.DB {
	return m.client.Master().Model(&User{}).Scopes(selectTable(uid))
}

func (m *Manager) slave(uid string) *gorm.DB {
	return m.client.Slave().Model(&User{}).Scopes(selectTable(uid))
}

func selectTable(uid string) func(tx *gorm.DB) *gorm.DB {
	if os.Getenv("PRETTYY_TEST") == "dev" {
		return func(tx *gorm.DB) *gorm.DB {
			return tx.Table(tablePrefix + "0")
		}
	}
	num := tool.Crc(uid, tableNum)
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Table(tablePrefix + num)
	}
}

func withUid(uid string) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Where("uid = ?", uid)
	}
}
