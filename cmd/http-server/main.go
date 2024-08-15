package main

import (
	"log"
	"os"
	"prettyy-server-online/cmd/http-server/api/article"
	"prettyy-server-online/cmd/http-server/api/base"
	"prettyy-server-online/cmd/http-server/api/common"
	"prettyy-server-online/cmd/http-server/api/user"
	"prettyy-server-online/cmd/http-server/auth"
	middleWare "prettyy-server-online/custom-pkg/xzf-gin-consul/middle-ware"
	ginConsulRegister "prettyy-server-online/custom-pkg/xzf-gin-consul/register"
)

func main() {
	container := ginConsulRegister.NewContainer()
	consulRegister := ginConsulRegister.NewRegisterConsul()
	if err := consulRegister.Init(); err != nil {
		log.Fatalf("new registry consul failed, err:%v\n", err)
	}
	container.SetRegistry(consulRegister)
	// 业务用一个端口，监控用另一个端口
	httpServer := ginConsulRegister.NewGinServer("blog-service")
	httpServer.SetAddress(os.Getenv("HTTP_SERVER_LISTEN_ADDR")) // 可支持通过环境变量来设置端口，如果不指定就是随机可用的端口号
	httpServer.Use(middleWare.Cors(), middleWare.NewZapLogger(), auth.Auth, middleWare.NewMetrics("blog_service_"))
	httpServer.AddService(
		base.NewServer(),
		user.NewServer(),
		article.NewServer(),
		common.NewServer())
	container.AddServer(httpServer)
	// 监控服务
	metricsServer := ginConsulRegister.NewGinServer("blog-service-metrics")
	metricsServer.SetAddress(os.Getenv("BLOG_METRICS_SERVER_LISTEN_ADDR")) // 可支持通过环境变量来设置端口，如果不指定就是随机可用的端口号
	metricsServer.AddService(middleWare.NewMetricService())
	container.AddServer(metricsServer)
	if err := container.Init(); err != nil {
		log.Fatalf("init service failed, err:%s\n", err.Error())
	}
	if err := container.Start(); err != nil {
		log.Fatalf("start fail: %s\n", err.Error())
	}
	if err := container.Wait(); err != nil {
		log.Fatalf("shutdown err: %s", err.Error())
	}
	return
}
