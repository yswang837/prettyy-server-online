package tool

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"hash/crc32"
	"path/filepath"
	"strconv"
	"time"
)

// Crc 分表的统一路由规则
func Crc(id string, tableNum int) string {
	return strconv.Itoa(int(crc32.ChecksumIEEE([]byte(id))) % tableNum)
}

func ToMd5(s string) string {
	b := md5.Sum([]byte(s))
	return hex.EncodeToString(b[:])
}

func Base64Encode(s string) string {
	if s == "" {
		return ""
	}
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func Base64Decode(s string) string {
	if s == "" {
		return ""
	}
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return ""
	}
	return string(b)
}

func StringToTime(s string) time.Time {
	if s == "" {
		return time.Time{}
	}
	t, err := time.Parse(DefaultDateTimeLayout, s)
	if err != nil {
		return time.Time{}
	}
	return t
}

func TimeToString(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(DefaultFileNameDateLayout)
}

func MakeFileName(filename string) string {
	if filename == "" {
		return ""
	}
	return TimeToString(time.Now()) + "_" + ToMd5(filename)[:8] + "." + filepath.Ext(filename)
}
