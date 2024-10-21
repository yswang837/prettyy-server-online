package register

import (
	"github.com/hashicorp/consul/api"
	"os"
	"strconv"
)

const (
	//定义tag信息
	tagMaintainer = "小钻风"
	tagVersion    = "1.0.1"
)

// RegisterConsul 定义consul的服务端的注册信息
type RegisterConsul struct {
	client *api.Client
	config ConsulConfig
}

// ConsulConfig 定义配置信息
type ConsulConfig struct {
	Address                 string
	CheckInterval           string
	CheckDeregisterInterval string
}

// NewRegisterConsul 创建一个consul的注册实例
func NewRegisterConsul() *RegisterConsul {
	return &RegisterConsul{
		config: ConsulConfig{
			CheckInterval:           "10s",
			CheckDeregisterInterval: "30s",
		},
	}
}

// Init 初始化consul的client信息及资源
func (r *RegisterConsul) Init() (err error) {
	conf := api.DefaultConfig()
	address := r.config.Address
	if address != "" {
		conf.Address = address
	}
	r.client, err = api.NewClient(conf)
	return
}

// Register 真正的注册 consul 服务实例
func (r *RegisterConsul) Register(serverName string, port int) error {
	check := &api.AgentServiceCheck{
		Interval:                       r.config.CheckInterval,
		TCP:                            r.getConsulServerIP() + ":" + strconv.Itoa(port),
		DeregisterCriticalServiceAfter: r.config.CheckDeregisterInterval,
	}
	registerOption := &api.AgentServiceRegistration{
		ID:    r.buildConsulServerID(serverName, port),
		Name:  serverName,
		Tags:  []string{tagMaintainer, tagVersion},
		Port:  port,
		Check: check,
	}
	return r.client.Agent().ServiceRegister(registerOption)
}

// DeRegister 注销 consul 服务实例
func (r *RegisterConsul) DeRegister(serverName string, port int) error {
	return r.client.Agent().ServiceDeregister(r.buildConsulServerID(serverName, port))
}

// getConsulServerIP 获取 consul 服务器IP
func (r *RegisterConsul) getConsulServerIP() string {
	ip := os.Getenv("CONSUL_SERVER_IP")
	if ip != "" {
		return ip
	}
	return "localhost"
}

// buildConsulServerID 构建 consul 服务实例ID
func (r *RegisterConsul) buildConsulServerID(serverName string, port int) string {
	return serverName + "_" + r.getConsulServerIP() + ":" + strconv.Itoa(port)
}

// SetAddress 支持设置地址
func (r *RegisterConsul) SetAddress(address string) {
	r.config.Address = address
}
