package xzf_qiniu

import (
	"context"
	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
	"mime/multipart"
)

func UploadFile(file multipart.File, fileSize int64) (string, error) {
	putPolicy := storage.PutPolicy{Scope: Bucket}
	mac := qbox.NewMac(AccessKey, SecretKey)
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{Zone: &storage.ZoneHuabei}
	putExtra := storage.PutExtra{}
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	if err := formUploader.PutWithoutKey(context.Background(), &ret, upToken, file, fileSize, &putExtra); err != nil {
		return "", err
	}
	return QiNiuServer + ret.Key, nil
}
