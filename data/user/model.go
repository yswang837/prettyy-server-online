package user

import (
	"encoding/json"
	"prettyy-server-online/utils/tool"
	"strconv"
	"time"
)

const (
	tableNum    = 2
	tablePrefix = "user_"
)

// User 面向数据库
type User struct {
	Uid           int64     `json:"uid"`            // 用户id，主键
	Email         string    `json:"email"`          // 邮箱，唯一索引
	Password      string    `json:"password"`       // 密码
	Phone         string    `json:"phone"`          // 电话号码
	NickName      string    `json:"nick_name"`      // 昵称
	Role          int       `json:"role"`           // 角色,用于权限管理,1普通用户,2管理员,3超级管理员
	Grade         int       `json:"grade"`          // 账号等级
	Avatar        string    `json:"avatar"`         // 头像
	Summary       string    `json:"summary"`        // 个人简介
	Gender        string    `json:"gender"`         // 男,女,保密
	ProvinceCity  string    `json:"province_city"`  // 户籍省市
	CodeAge       int       `json:"code_age"`       // 码龄
	IsCertified   int       `json:"is_certified"`   // 是否认证, 0未认证, 1认证，邮箱和电话都绑定即为认证
	DataIntegrity int       `json:"data_integrity"` // 资料完整度
	Birthday      string    `json:"birthday"`       // 生日格式 xxxx-xx-xx
	CreateTime    time.Time `json:"create_time"`    // 创建时间
	UpdateTime    time.Time `json:"update_time"`    // 更新时间
	LoginTime     time.Time `json:"login_time"`     // 登录时间
}

func (u *User) TableName() string {
	return tablePrefix + tool.Crc(strconv.FormatInt(u.Uid, 10), tableNum)
}

func (u *User) String() string {
	b, _ := json.Marshal(u)
	return string(b)
}
