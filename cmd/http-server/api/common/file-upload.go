package common

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	ginConsulRegister "prettyy-server-online/custom-pkg/xzf-gin-consul/register"
	"prettyy-server-online/utils/metrics"
	"prettyy-server-online/utils/tool"
)

// FileUpload 文件上传
// 4000140
// 2000140

func (s *Server) FileUpload(ctx *ginConsulRegister.Context) {
	metrics.CommonCounter.Inc("file-upload", "total")
	file, err := ctx.FormFile("file")
	if err != nil {
		metrics.CommonCounter.Inc("file-upload", "params-error")
		ctx.SetError(err.Error())
		ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000140, Message: "参数错误"})
		return
	}
	src, err := file.Open()
	if err != nil {
		metrics.CommonCounter.Inc("file-upload", "open-error")
		ctx.SetError(err.Error())
		ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000141, Message: "打开文件错误"})
		return
	}
	defer src.Close()
	filename := tool.MakeFileName(file.Filename)
	ctx.SetFilename(filename)
	dst, err := os.Create(filepath.Join(uploadDir, filename))
	if err != nil {
		metrics.CommonCounter.Inc("file-upload", "touch-error")
		ctx.SetError(err.Error())
		ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000142, Message: "创建文件错误"})
		return
	}
	defer dst.Close()
	if _, err = io.Copy(dst, src); err != nil {
		metrics.CommonCounter.Inc("file-upload", "save-error")
		ctx.SetError(err.Error())
		ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000143, Message: "保存文件错误"})
		return
	}
	metrics.CommonCounter.Inc("file-upload", "succ")
	if os.Getenv("idc") == "dev" {
		// 测试时，需关闭auth中间件
		ctx.JSON(http.StatusOK, ginConsulRegister.Response{Code: 2000140, Message: "上传文件成功", Result: fmt.Sprintf("http://127.0.0.1:6677/uploads/%s", tool.MakeFileName(file.Filename))})
		return
	}
	ctx.JSON(http.StatusOK, ginConsulRegister.Response{Code: 2000140, Message: "上传文件成功", Result: fmt.Sprintf("http://120.26.203.121/uploads/%s", tool.MakeFileName(file.Filename))})
	return
}
