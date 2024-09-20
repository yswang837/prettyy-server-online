package common

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"path/filepath"
	ginConsulRegister "prettyy-server-online/custom-pkg/xzf-gin-consul/register"
)

// FileUpload 文件上传
// 4000140
// 2000140

func (s *Server) FileUpload(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000140, Message: "参数错误"})
		return
	}
	src, err := file.Open()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000141, Message: "打开文件错误"})
		return
	}
	defer src.Close()
	dst, err := os.Create(filepath.Join(uploadDir, file.Filename))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000142, Message: "创建文件错误"})
		return
	}
	defer dst.Close()
	if _, err = io.Copy(dst, src); err != nil {
		ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000143, Message: "保存文件错误"})
		return
	}
	ctx.JSON(http.StatusOK, ginConsulRegister.Response{Code: 2000140, Message: "上传文件成功", Result: fmt.Sprintf("http://120.26.203.121/uploads/%s", file.Filename)})
	return
}
