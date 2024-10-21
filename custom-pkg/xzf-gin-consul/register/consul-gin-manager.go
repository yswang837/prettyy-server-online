package register

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Container 定义consul和gin-server的Container，主函数启动本Container，就可以直接启动gin-server + consul
type Container struct {
	registry *RegisterConsul
	servers  map[string]ManagerServer
	stopCh   chan os.Signal
}

// ManagerServer 定义了一个服务的接口，包括服务的启动、关闭、注册、服务地址、服务端口等，GinServer实现了该接口
type ManagerServer interface {
	Name() string
	Init() error
	StartListen() error
	Address() string
	Port() int
	Shutdown(ctx context.Context) error
}

// NewContainer 创建一个Container
func NewContainer() *Container {
	return &Container{servers: map[string]ManagerServer{}, stopCh: make(chan os.Signal)}
}

// AddServer 向Container中添加实现了ManagerServer接口的所有服务
func (c *Container) AddServer(servers ...ManagerServer) {
	for _, server := range servers {
		c.servers[server.Name()] = server
	}
}

// Init 初始化Container中实现了ManagerServer接口的所有服务
func (c *Container) Init() error {
	for _, server := range c.servers {
		if err := server.Init(); err != nil {
			return err
		}
	}
	return nil
}

// Start 启动Container中实现了ManagerServer接口的所有服务
func (c *Container) Start() error {
	ch := make(chan error, len(c.servers))
	for _, server := range c.servers {
		// 异步启动所有服务
		go func(server ManagerServer) {
			if err := server.StartListen(); err != nil {
				ch <- err
			}
		}(server)
	}
	// select多路复用，阻塞等待2秒start所有服务，2秒后没有出错则认为所有服务都start成功
	var err error
	select {
	case <-time.After(2 * time.Second):
	case err = <-ch:
	}
	if err != nil {
		return err
	}
	// 服务start成功，将所有服务注册给consul，由consul统一接管

	// 对于每一个实现了ManagerServer接口的服务，获取它的ip和端口绑定到server上，供后续使用
	err = c.ergodicService(func(server ManagerServer) error {
		address := server.Address()
		port := server.Port()
		if port == 0 {
			return fmt.Errorf("%s's address invalid", address)
		}
		return nil
	})
	if err != nil {
		return err
	}

	// 对于每一个实现了ManagerServer接口的服务，将其注册到consul
	err = c.ergodicService(func(server ManagerServer) error {
		if err = c.registry.Register(server.Name(), server.Port()); err != nil {
			return fmt.Errorf("register %s %s fail: %s", server.Name(), server.Address(), err.Error())
		}
		log.Printf("success register %s on %s", server.Name(), server.Address())
		return nil
	})

	return err
}

// Shutdown 注销Container中实现了ManagerServer接口的所有服务，从consul上抹掉，并且Shutdown所有服务
func (c *Container) shutdown(context context.Context) error {
	_ = c.ergodicService(func(server ManagerServer) error {
		if err := c.registry.DeRegister(server.Name(), server.Port()); err != nil {
			log.Printf("deregister %s(%s) fail: %s", server.Name(), server.Address(), err.Error())
			return err
		}
		log.Printf("deregister %s(%s) succ", server.Name(), server.Address())
		return nil
	})
	_ = c.ergodicService(func(server ManagerServer) error {
		log.Printf("shutdown %s ...", server.Name())
		return server.Shutdown(context)
	})
	signal.Stop(c.stopCh)
	close(c.stopCh)
	return nil
}

// ergodicService 每一个server传递给f干啥由调用方指定
func (c *Container) ergodicService(f func(server ManagerServer) error) (err error) {
	for _, server := range c.servers {
		if err = f(server); err != nil {
			break
		}
	}
	return err
}

// SetRegistry 需要时单独初始化(未放在NewContainer的参数里面) consul注册，Container只启动gin-server也行
func (c *Container) SetRegistry(registry *RegisterConsul) {
	c.registry = registry
}

// Wait 等待退出信号，并且退出
func (c *Container) Wait() error {
	var err error

	signal.Notify(c.stopCh, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	sig := <-c.stopCh
	log.Println("stop by ", sig)
	chWait := make(chan error)
	// 异步shutdown
	go func() {
		chWait <- c.shutdown(context.Background())
	}()
	select {
	case err = <-chWait:
	case <-time.After(5 * time.Second):
		err = errors.New("shutdown timeout")
	}
	return err
}
