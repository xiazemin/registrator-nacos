package nacos

import (
	"log"
	"net/url"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/xiazemin/registrator/bridge"
)

func init() {
	bridge.Register(new(Factory), "nacos")
}

type Factory struct{}

func (f *Factory) New(uri *url.URL) bridge.RegistryAdapter {
	serverConfig := make([]constant.ServerConfig, 0)
	// 创建serverConfig
	// 支持多个;至少一个ServerConfig
	if uri.Host != "" {
		serverConfig = append(serverConfig, constant.ServerConfig{
			IpAddr: uri.Host,
			Port:   8848,
		})
	} else {
		serverConfig = append(serverConfig, constant.ServerConfig{
			IpAddr: "127.0.0.1",
			Port:   8848,
		})
	}

	// 创建clientConfig
	clientConfig := constant.ClientConfig{
		NamespaceId:         "", // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           50000,
		NotLoadCacheAtStart: true,
		LogLevel:            "debug",
	}

	// 创建服务发现客户端的另一种方式 (推荐)
	namingClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfig,
		},
	)
	if err != nil {
		log.Fatalf("初始化nacos失败: %s", err.Error())
	}

	return &NacosAdapter{client: namingClient, path: uri.Path}
}

type NacosAdapter struct {
	client naming_client.INamingClient
	path   string
}

func (r *NacosAdapter) Ping() error {
	var err error
	if r.client != nil {
		_, err = r.client.GetAllServicesInfo(vo.GetAllServiceInfoParam{
			GroupName: "DEFAULT_GROUP", // 默认值DEFAULT_GROUP
			NameSpace: "DEFAULT",       // 默认值DEFAULT
			PageNo:    0,
			PageSize:  1,
		})
	}

	if err != nil {
		return err
	}
	return nil
}

func (r *NacosAdapter) Register(service *bridge.Service) error {
	if service.IP == "" {
		log.Println(service)
		return nil
	}
	success, err := r.client.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          service.IP,
		Port:        uint64(service.Port),
		ServiceName: service.Name,
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    map[string]string{"name": service.Name, "id": service.ID},
		ClusterName: "DEFAULT",       // 默认值DEFAULT
		GroupName:   "DEFAULT_GROUP", // 默认值DEFAULT_GROUP
	})
	if err != nil || !success {
		log.Fatalf("注册服务失败: %s ,ip:%s,port:%d,name:%s", err.Error(), service.IP, service.Port, service.Name)
	}
	return err
}

func (r *NacosAdapter) Deregister(service *bridge.Service) error {
	success, err := r.client.DeregisterInstance(vo.DeregisterInstanceParam{
		Ip:          service.IP,
		Port:        uint64(service.Port),
		ServiceName: service.Name,
		Cluster:     "DEFAULT",
		GroupName:   "DEFAULT_GROUP", // 默认值DEFAULT_GROUP
		Ephemeral:   true,
	})
	if err != nil || !success {
		log.Fatalf("取消注册服务失败: %s", err.Error())
	}
	return err
}

func (r *NacosAdapter) Refresh(service *bridge.Service) error {
	return r.Register(service)
}

func (r *NacosAdapter) Services() ([]*bridge.Service, error) {
	return []*bridge.Service{}, nil
}
