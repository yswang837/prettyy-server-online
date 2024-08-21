package xzf_email

import (
	"fmt"
	"github.com/jordan-wright/email"
	"math/rand"
	"net/smtp"
	"time"
)

func SendEmail(to string) (string, error) {
	e := email.NewEmail()
	e.From = EmailFrom
	e.Subject = Subject
	e.To = []string{to}
	code := fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
	str := "【验证码】为：" + code + ", 仅" + ActiveTime + "分钟内有效,请勿转发他人!"
	e.HTML = []byte(str)
	if err := e.Send(Addr+Port, smtp.PlainAuth("", Username, Password, Addr)); err != nil {
		return "", err
	}
	return code, nil
}
