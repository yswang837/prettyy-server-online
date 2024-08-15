package xzf_email

import (
	"fmt"
	"gopkg.in/ini.v1"
	"log"
	"os"
)

var (
	EmailFrom  string
	Addr       string
	Port       string
	Username   string
	Password   string
	ActiveTime string
	Subject    string
)

func init() {
	pwd, _ := os.Getwd()
	fmt.Println("pwd.....", pwd)
	f, err := ini.Load(pwd + "/config/config.ini")
	if err != nil {
		log.Fatalf("read config file err: %s", err.Error())
	}
	LoadEmail(f)
}

func LoadEmail(f *ini.File) {
	EmailFrom = f.Section("email").Key("EmailFrom").MustString("小钻风科技有限责任公司 <1714113169@qq.com>")
	Addr = f.Section("email").Key("Addr").MustString("smtp.qq.com")
	Port = f.Section("email").Key("Port").MustString(":587")
	Username = f.Section("email").Key("Username").MustString("1714113169@qq.com")
	Password = f.Section("email").Key("Password").MustString("bxvjwuyyeoqfdced")
	ActiveTime = f.Section("email").Key("ActiveTime").MustString("5")
	Subject = f.Section("email").Key("Subject").MustString("验证码")
}
