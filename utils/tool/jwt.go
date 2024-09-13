package tool

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"strings"
	"time"
)

const (
	TokenExpire = 86400
)

// jwt秘钥
var stSignKey = []byte("cab670bd4fcc32833476660de1ad1056")

// JwtCustomClaims 注册声明是JWT声明集的结构化版本，仅限于注册声明名称，先把uid属性删掉，后续需要再还原
type JwtCustomClaims struct {
	Uid              int64
	RegisteredClaims jwt.RegisteredClaims
}

func (j JwtCustomClaims) Valid() error {
	return nil
}

// GenerateToken 生成Token
func GenerateToken(uid int64) (string, error) {
	// 初始化
	iJwtCustomClaims := JwtCustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			// 设置过期时间 在当前基础上 添加一天后 过期
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(TokenExpire) * time.Second)),
			// 颁发时间 也就是生成时间
			IssuedAt: jwt.NewNumericDate(time.Now()),
			//主题
			Subject: "Token",
		},
		Uid: uid,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, iJwtCustomClaims)
	return token.SignedString(stSignKey)
}

// ParseToken 解析token
func ParseToken(tokenStr string) (JwtCustomClaims, error) {
	iJwtCustomClaims := JwtCustomClaims{}
	//ParseWithClaims是NewParser().ParseWithClaims()的快捷方式
	token, err := jwt.ParseWithClaims(tokenStr, &iJwtCustomClaims, func(token *jwt.Token) (interface{}, error) {
		return stSignKey, nil
	})

	if err == nil && !token.Valid {
		err = errors.New("invalid Token")
	}
	return iJwtCustomClaims, err
}

func IsTokenValid(tokenStr string) bool {
	_, err := ParseToken(tokenStr)
	fmt.Println(err)
	if err != nil {
		return false
	}
	return true
}

func GetToken(c *gin.Context) string {
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		return ""
	}
	splitToken := strings.Split(token, " ")
	if len(splitToken) <= 1 {
		return ""
	}
	return splitToken[1]
}
