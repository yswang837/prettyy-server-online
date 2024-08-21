package user

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
	if err := m.master(u.Email).Create(u).Error; err != nil {
		return err
	}
	return nil
}

func (m *Manager) Get(email string) (*User, error) {
	if email == "" {
		return nil, tool.ErrParams
	}
	u := &User{}
	if err := m.slave(email).Scopes(withEmail(email)).First(u).Error; err != nil {
		return nil, err
	}
	if u.Email == "" {
		return nil, errors.New("record not found")
	}
	return u, nil
}

func (m *Manager) UpdateLoginTime(email string) error {
	if email == "" {
		return tool.ErrParams
	}
	if err := m.master(email).Scopes(withEmail(email)).Update("login_time", time.Now()).Error; err != nil {
		return err
	}
	return nil
}

func (m *Manager) UpdatePassword(email, password string) error {
	if email == "" || password == "" {
		return tool.ErrParams
	}
	if err := m.master(email).Scopes(withEmail(email)).Update("password", password).Error; err != nil {
		return err
	}
	return nil
}

func (m *Manager) UpdateNickName(email, nickName string) error {
	if email == "" || nickName == "" {
		return tool.ErrParams
	}
	if err := m.master(email).Scopes(withEmail(email)).Update("nick_name", nickName).Error; err != nil {
		return err
	}
	return nil
}

func (m *Manager) UpdateSummary(email, summary string) error {
	if email == "" || summary == "" {
		return tool.ErrParams
	}
	if err := m.master(email).Scopes(withEmail(email)).Update("summary", summary).Error; err != nil {
		return err
	}
	return nil
}

func (m *Manager) UpdateProvinceCity(email, provinceCity string) error {
	if email == "" || provinceCity == "" {
		return tool.ErrParams
	}
	if err := m.master(email).Scopes(withEmail(email)).Update("province_city", provinceCity).Error; err != nil {
		return err
	}
	return nil
}

func (m *Manager) UpdateBirthdayCity(email, birthday string) error {
	if email == "" || birthday == "" {
		return tool.ErrParams
	}
	if err := m.master(email).Scopes(withEmail(email)).Update("birthday", birthday).Error; err != nil {
		return err
	}
	return nil
}

func (m *Manager) UpdateGender(email, gender string) error {
	if email == "" || gender == "保密" || gender == "" {
		return tool.ErrParams
	}
	if err := m.master(email).Scopes(withEmail(email)).Update("gender", gender).Error; err != nil {
		return err
	}
	return nil
}

func (m *Manager) master(email string) *gorm.DB {
	return m.client.Master().Model(&User{}).Scopes(selectTable(email))
}

func (m *Manager) slave(email string) *gorm.DB {
	return m.client.Slave().Model(&User{}).Scopes(selectTable(email))
}

func selectTable(tid string) func(tx *gorm.DB) *gorm.DB {
	if os.Getenv("PRETTYY_TEST") == "dev" {
		return func(tx *gorm.DB) *gorm.DB {
			return tx.Table(tablePrefix + "0")
		}
	}
	num := tool.Crc(tid, tableNum)
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Table(tablePrefix + num)
	}
}

func withEmail(email string) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Where("email = ?", email)
	}
}
