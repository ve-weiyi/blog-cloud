package nacos

import (
	"log"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

type NacosConfig struct {
	IP          string
	Port        uint64 //8848
	UserName    string
	Password    string
	NameSpaceID string
	Group       string
	DataID      string
	LogLevel    string //debug
	Timeout     int64  //ms
}

type NacosReader struct {
	cfg *NacosConfig
}

func New(cfg *NacosConfig) *NacosReader {
	return &NacosReader{cfg: cfg}
}

func (n *NacosReader) Init(listener func(content string) error) error {
	var dataId, group string
	dataId = n.cfg.DataID
	group = n.cfg.Group

	//create ServerConfig
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(n.cfg.IP, n.cfg.Port),
	}

	//create ClientConfig
	cc := constant.NewClientConfig(
		constant.WithUsername(n.cfg.UserName),
		constant.WithPassword(n.cfg.Password),
		constant.WithNamespaceId(n.cfg.NameSpaceID),
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir("./runtime/logs"),
		constant.WithCacheDir("./runtime/cache"),
		constant.WithLogLevel("debug"),
	)

	// create config client
	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		return err
	}

	//get config
	content, err := client.GetConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
	})
	if err != nil {
		return err
	}

	log.Println("nacos get config :"+content, err)

	err = listener(content)
	if err != nil {
		return err
	}

	//Listen config change,key=dataId+group+namespaceId.
	err = client.ListenConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
		OnChange: func(namespace, group, dataId, data string) {
			log.Println("nacos config changed group:" + group + ", dataId:" + dataId + ", content:" + data)
			if err = listener(data); err != nil {
				log.Println("nacos config changed reload failed")
			}
		},
	})
	if err != nil {
		return err
	}

	return nil
}
