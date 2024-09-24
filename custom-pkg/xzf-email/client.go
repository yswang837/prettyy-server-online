package xzf_email

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"math/rand"
	"time"
)

func SendEmail(to string) (string, error) {
	e := gomail.NewMessage()
	e.SetHeader("From", e.FormatAddress(EmailFrom, EmailName))
	e.SetHeader("Subject", Subject)
	e.SetHeader("To", to)
	code := fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
	str := "【验证码】为：" + code + ", 仅" + ActiveTime + "分钟内有效,请勿转发他人!"
	e.SetBody("text/html", str)
	if err := gomail.NewDialer(Addr, Port, Username, Password).DialAndSend(e); err != nil {
		return "", err
	}
	return code, nil
}
