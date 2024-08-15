package register

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

//gin的路由及参数获取等+go自带的http.server

// GinServer 实现了ManagerServer接口，ManagerServer接口的变量可直接调用绑定在GinServer上面的方法
type GinServer struct {
	name     string
	address  string
	r        *gin.Engine
	server   *http.Server
	services []HttpServer
}

type HttpServer interface {
	Init() error
	SetRoute(r *gin.Engine) //绑定路由和具体的函数执行逻辑
}

func NewGinServer(name string) *GinServer {
	setGinMode()
	r := gin.New()
	r.Use(gin.Recovery())
	return &GinServer{r: r, name: name}
}

func (g *GinServer) SetGroup(groupName string) *GinServer {
	g2 := g
	g2.r.Group(groupName)
	return g2
}

func setGinMode() {
	mode := os.Getenv("GIN_MODE")
	switch mode {
	case "debug":
		gin.SetMode(gin.DebugMode)
	case "release":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}
}

func (g *GinServer) Name() string {
	return g.name
}

func (g *GinServer) Init() error {
	server := &http.Server{Handler: g.r.Handler(), IdleTimeout: 10 * time.Second}
	g.server = server

	for _, service := range g.services {
		service.SetRoute(g.r)
		if err := service.Init(); err != nil {
			return err
		}
	}
	return nil
}

func (g *GinServer) AddService(services ...HttpServer) {
	for _, service := range services {
		g.services = append(g.services, service)
	}
}

func (g *GinServer) Use(middleWare ...gin.HandlerFunc) {
	g.r.Use(middleWare...)
}

func (g *GinServer) StartListen() error {
	address := ":0"
	if g.address != "" {
		address = g.address
	}
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("listen failed, err:%v\n", err)
	}
	g.SetAddress(listener.Addr().String())
	if err = g.server.Serve(listener); err != nil {
		return fmt.Errorf("startup service failed, err:%v\n", err)
	}
	return nil
}

func (g *GinServer) Address() string {
	return g.address
}

func (g *GinServer) Port() int {
	_, portString, _ := net.SplitHostPort(g.address)
	port, _ := strconv.Atoi(portString)
	return port
}

func (g *GinServer) SetAddress(address string) {
	g.address = address
}

func (g *GinServer) Shutdown(ctx context.Context) error {
	return g.server.Shutdown(ctx)
}
