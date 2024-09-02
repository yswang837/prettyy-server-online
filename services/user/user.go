package user

import (
	"errors"
	"prettyy-server-online/data/user"
	"prettyy-server-online/utils/tool"
	"strconv"
	"strings"
	"time"
)

const (
	role               = 1   // 通过接口只允许注册普通用户，若注册管理员，则用命令行工具去update
	grade              = 1   // 注册时账号等级为1
	codeAge            = 1   // 码龄1天
	dataIntegrity      = 30  // 注册用户，资料完整度默认30%
	identifyCodeExpire = 300 // 验证码有效期5分钟
	avatar             = "https://s21.ax1x.com/2024/07/05/pkRgyT0.jpg"
	uidIncrKey         = "bj-blog-uid"
	gender             = "保密" // 性别未知
)

type Client struct {
	manager      *user.Manager
	cacheManager *user.ManagerRedis
}

var defaultClient *Client

func init() {
	var err error
	defaultClient, err = NewClient()
	if err != nil {
		panic(err)
	}
}

func NewClient() (*Client, error) {
	manager, err := user.NewManager()
	if err != nil {
		return nil, err
	}
	cache, err := user.NewManagerRedis()
	if err != nil {
		return nil, err
	}
	return &Client{manager: manager, cacheManager: cache}, nil
}

func Add(u *user.User) (err error) {
	return defaultClient.Add(u)
}

func (c *Client) Add(u *user.User) (err error) {
	if u == nil {
		return tool.ErrParams
	}
	u.Uid, err = c.Incr(uidIncrKey)
	if err != nil {
		return errors.New("generate uid err: " + err.Error())
	}
	if u.Password != "" {
		u.Password = tool.ToMd5(u.Password)
	}
	u.Avatar = avatar
	u.Role = role
	u.Grade = grade
	u.CodeAge = codeAge
	u.Gender = gender
	u.DataIntegrity = dataIntegrity
	u.NickName = strings.Split(u.Email, "@")[0] // 默认用户名用邮箱前缀代替
	if err = c.manager.Add(u); err != nil {
		return errors.New("register to mysql failed: " + err.Error())
	}
	u.CreateTime = time.Now()
	u.UpdateTime = time.Now()
	u.LoginTime = time.Now()
	if _, err = c.cacheManager.HMSet(strconv.FormatInt(u.Uid, 10), userToMap(u)); err != nil {
		return errors.New("register to redis failed: " + err.Error())
	}
	return
}

func userToMap(u *user.User) map[string]interface{} {
	if u == nil {
		return nil
	}
	m := make(map[string]interface{})
	m["uid"] = u.Uid
	m["email"] = u.Email
	m["password"] = u.Password
	m["phone"] = u.Phone
	m["nick_name"] = u.NickName
	m["role"] = u.Role
	m["grade"] = u.Grade
	m["avatar"] = u.Avatar
	m["summary"] = u.Summary
	m["gender"] = u.Gender
	m["code_age"] = u.CodeAge
	m["is_certified"] = u.IsCertified
	m["data_integrity"] = u.DataIntegrity
	m["create_time"] = u.CreateTime.Format(tool.DefaultDateTimeLayout)
	m["update_time"] = u.UpdateTime.Format(tool.DefaultDateTimeLayout)
	m["login_time"] = u.LoginTime.Format(tool.DefaultDateTimeLayout)
	return m
}

func mapToUser(m map[string]string) *user.User {
	if m == nil {
		return nil
	}
	r, _ := strconv.Atoi(m["role"])
	g, _ := strconv.Atoi(m["grade"])
	c, _ := strconv.Atoi(m["code_age"])
	cer, _ := strconv.Atoi(m["is_certified"])
	uid, _ := strconv.Atoi(m["uid"])

	u := &user.User{
		Uid:         int64(uid),
		Email:       m["email"],
		Password:    m["password"],
		Phone:       m["phone"],
		NickName:    m["nick_name"],
		Avatar:      m["avatar"],
		Role:        r,
		Grade:       g,
		Summary:     m["summary"],
		Gender:      m["gender"],
		CodeAge:     c,
		IsCertified: cer,
		CreateTime:  tool.StringToTime(m["create_time"]),
		UpdateTime:  tool.StringToTime(m["update_time"]),
		LoginTime:   tool.StringToTime(m["login_time"]),
	}
	return u
}

func (c *Client) GetUser(uid string) (*user.User, error) {
	if uid == "" {
		return nil, tool.ErrParams
	}
	m, err := c.cacheManager.HGetAll(uid)
	if err != nil {
		return nil, err
	}
	var u *user.User
	if len(m) != 0 {
		u = mapToUser(m)
		if u != nil {
			return u, nil
		}
	}
	u, err = c.manager.Get(uid)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func UpdateLoginTime(uid string) error {
	return defaultClient.UpdateLoginTime(uid)
}

func (c *Client) UpdateLoginTime(uid string) error {
	if uid == "" {
		return tool.ErrParams
	}
	_, err := c.cacheManager.HSet(uid, "login_time", time.Now().Format(tool.DefaultDateTimeLayout))
	if err != nil {
		return errors.New("redis update login time failed: " + err.Error())
	}
	return c.manager.UpdateLoginTime(uid)
}

func UpdatePassword(uid, password string) error {
	return defaultClient.UpdatePassword(uid, password)
}

func (c *Client) UpdatePassword(uid, password string) error {
	if uid == "" || password == "" {
		return tool.ErrParams
	}
	password = tool.ToMd5(password)
	_, err := c.cacheManager.HSet(uid, "password", password)
	if err != nil {
		return errors.New("redis update password failed: " + err.Error())
	}
	return c.manager.UpdatePassword(uid, password)
}

func UpdateNickName(uid, nickName string) error {
	return defaultClient.UpdateNickName(uid, nickName)
}

func (c *Client) UpdateNickName(uid, nickName string) error {
	if uid == "" || nickName == "" {
		return tool.ErrParams
	}
	_, err := c.cacheManager.HSet(uid, "nick_name", nickName)
	if err != nil {
		return errors.New("redis update nick name failed: " + err.Error())
	}
	return c.manager.UpdateNickName(uid, nickName)
}

func UpdateSummary(uid, summary string) error {
	return defaultClient.UpdateSummary(uid, summary)
}

func (c *Client) UpdateSummary(uid, summary string) error {
	if uid == "" || summary == "" {
		return tool.ErrParams
	}
	_, err := c.cacheManager.HSet(uid, "summary", summary)
	if err != nil {
		return errors.New("redis update summary failed: " + err.Error())
	}
	return c.manager.UpdateSummary(uid, summary)
}

func UpdateProvinceCity(uid, provinceCity string) error {
	return defaultClient.UpdateProvinceCity(uid, provinceCity)
}

func (c *Client) UpdateProvinceCity(uid, UpdateProvinceCity string) error {
	if uid == "" || UpdateProvinceCity == "" {
		return tool.ErrParams
	}
	_, err := c.cacheManager.HSet(uid, "province_city", UpdateProvinceCity)
	if err != nil {
		return errors.New("redis update province_city failed: " + err.Error())
	}
	return c.manager.UpdateProvinceCity(uid, UpdateProvinceCity)
}

func UpdateBirthdayCity(uid, birthday string) error {
	return defaultClient.UpdateBirthdayCity(uid, birthday)
}

func (c *Client) UpdateBirthdayCity(uid, birthday string) error {
	if uid == "" || birthday == "" {
		return tool.ErrParams
	}
	_, err := c.cacheManager.HSet(uid, "birthday", birthday)
	if err != nil {
		return errors.New("redis update birthday failed: " + err.Error())
	}
	return c.manager.UpdateBirthdayCity(uid, birthday)
}

func UpdateGender(uid, gender string) error {
	return defaultClient.UpdateGender(uid, gender)
}

func (c *Client) UpdateGender(uid, g string) error {
	if uid == "" || g == gender || g == "" {
		return tool.ErrParams
	}
	u, err := c.GetUser(uid)
	if err != nil {
		return errors.New("update gender get user failed: " + err.Error())
	}
	if u.Gender == "男" || u.Gender == "女" {
		return errors.New("can not change gender")
	}
	_, err = c.cacheManager.HSet(uid, "gender", g)
	if err != nil {
		return errors.New("redis update gender failed: " + err.Error())
	}
	return c.manager.UpdateGender(uid, g)
}

func GetUser(uid string) (*user.User, error) {
	return defaultClient.GetUser(uid)
}

func (c *Client) SetEx(uid string, value string, expire int) error {
	_, err := c.cacheManager.SetWithTTL(uid, expire, value)
	return err
}

func (c *Client) Incr(key string) (int64, error) {
	return c.cacheManager.Incr(key)
}

// SetExByEmail set邮箱验证码并设置5分钟过期
func SetExByEmail(email, value string) error {
	return defaultClient.SetEx(buildIdentifyCode(email), value, identifyCodeExpire)
}

func SetExByToken(token string) error {
	return defaultClient.SetEx(token, "1", tool.TokenExpire)
}

func IsExistToken(token string) bool {
	return defaultClient.IsExistToken(token)
}

func (c *Client) IsExistToken(token string) bool {
	val, _ := c.cacheManager.Get(token)
	if val == "1" {
		return true
	}
	return false
}

func (c *Client) GetIdentifyCodeFromCache(email string) string {
	s, _ := c.cacheManager.Get(email)
	return s
}

// GetIdentifyCodeFromCache 通过email获取验证码
func GetIdentifyCodeFromCache(email string) string {
	return defaultClient.GetIdentifyCodeFromCache(buildIdentifyCode(email))
}

func buildIdentifyCode(email string) string {
	return email + ":code"
}
