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
	if _, err = c.cacheManager.HMSet(u.Email, userToMap(u)); err != nil {
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

func Add(u *user.User) (err error) {
	return defaultClient.Add(u)
}

func (c *Client) GetUser(email string) (*user.User, error) {
	if email == "" {
		return nil, tool.ErrParams
	}
	m, err := c.cacheManager.HGetAll(email)
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
	u, err = c.manager.Get(email)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (c *Client) UpdateLoginTime(email string) error {
	if email == "" {
		return tool.ErrParams
	}
	_, err := c.cacheManager.HSet(email, "login_time", time.Now().Format(tool.DefaultDateTimeLayout))
	if err != nil {
		return errors.New("redis update login time failed: " + err.Error())
	}
	return c.manager.UpdateLoginTime(email)
}

func UpdateLoginTime(email string) error {
	return defaultClient.UpdateLoginTime(email)
}

func (c *Client) UpdatePassword(email, password string) error {
	if email == "" || password == "" {
		return tool.ErrParams
	}
	_, err := c.cacheManager.HSet(email, "password", password)
	if err != nil {
		return errors.New("redis update password failed: " + err.Error())
	}
	return c.manager.UpdatePassword(email, password)
}

func UpdatePassword(email, password string) error {
	return defaultClient.UpdatePassword(email, password)
}

func UpdateNickName(email, nickName string) error {
	return defaultClient.UpdateNickName(email, nickName)
}

func (c *Client) UpdateNickName(email, nickName string) error {
	if email == "" || nickName == "" {
		return tool.ErrParams
	}
	_, err := c.cacheManager.HSet(email, "nick_name", nickName)
	if err != nil {
		return errors.New("redis update nick name failed: " + err.Error())
	}
	return c.manager.UpdateNickName(email, nickName)
}

func UpdateSummary(email, summary string) error {
	return defaultClient.UpdateSummary(email, summary)
}

func UpdateProvinceCity(email, provinceCity string) error {
	return defaultClient.UpdateProvinceCity(email, provinceCity)
}

func UpdateBirthdayCity(email, birthday string) error {
	return defaultClient.UpdateBirthdayCity(email, birthday)
}

func (c *Client) UpdateBirthdayCity(email, birthday string) error {
	if email == "" || birthday == "" {
		return tool.ErrParams
	}
	_, err := c.cacheManager.HSet(email, "birthday", birthday)
	if err != nil {
		return errors.New("redis update birthday failed: " + err.Error())
	}
	return c.manager.UpdateBirthdayCity(email, birthday)
}

func (c *Client) UpdateProvinceCity(email, UpdateProvinceCity string) error {
	if email == "" || UpdateProvinceCity == "" {
		return tool.ErrParams
	}
	_, err := c.cacheManager.HSet(email, "province_city", UpdateProvinceCity)
	if err != nil {
		return errors.New("redis update province_city failed: " + err.Error())
	}
	return c.manager.UpdateProvinceCity(email, UpdateProvinceCity)
}

func (c *Client) UpdateSummary(email, summary string) error {
	if email == "" || summary == "" {
		return tool.ErrParams
	}
	_, err := c.cacheManager.HSet(email, "summary", summary)
	if err != nil {
		return errors.New("redis update summary failed: " + err.Error())
	}
	return c.manager.UpdateSummary(email, summary)
}

func UpdateGender(email, gender string) error {
	return defaultClient.UpdateGender(email, gender)
}

func (c *Client) UpdateGender(email, g string) error {
	if email == "" || g == gender || g == "" {
		return tool.ErrParams
	}
	u, err := c.GetUser(email)
	if err != nil {
		return errors.New("update gender get user failed: " + err.Error())
	}
	if u.Gender == "男" || u.Gender == "女" {
		return errors.New("can not change gender")
	}
	_, err = c.cacheManager.HSet(email, "gender", g)
	if err != nil {
		return errors.New("redis update gender failed: " + err.Error())
	}
	return c.manager.UpdateGender(email, g)
}

func GetUser(email string) (*user.User, error) {
	return defaultClient.GetUser(email)
}

func (c *Client) SetEx(email string, value string, expire int) error {
	_, err := c.cacheManager.SetWithTTL(email, expire, value)
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

func init() {
	var err error
	defaultClient, err = NewClient()
	if err != nil {
		panic(err)
	}
	//已改为用redis自增key的方式来生成uid
	//if err = xzfSnowflake.Init("2024-03-09", "1"); err != nil {
	//	panic(err)
	//}
}
