package xzf_qiniu

import (
	"gopkg.in/ini.v1"
	"log"
	"os"
)

var (
	AccessKey   string
	SecretKey   string
	Bucket      string
	QiNiuServer string
)

func init() {
	pwd, _ := os.Getwd()
	f, err := ini.Load(pwd + "/config/qiniu/default.ini")
	if err != nil {
		log.Fatalf("read config file err: %s", err.Error())
	}
	LoadQiNiu(f)

}

func LoadQiNiu(f *ini.File) {
	AccessKey = f.Section("qiniu").Key("AccessKey").String()
	SecretKey = f.Section("qiniu").Key("SecretKey").String()
	Bucket = f.Section("qiniu").Key("Bucket").String()
	QiNiuServer = f.Section("qiniu").Key("QiNiuServer").String()
}
