package nacos

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/common/logger"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

type NacosConfig struct {
	IP          string
	Port        uint64 //8848
	UserName    string
	Password    string
	NameSpaceId string
	Group       string
	DataId      string
	RuntimeDir  string //runtime
	LogLevel    string //debug
	Timeout     int64  //ms
}

type NacosReader struct {
	cfg    *NacosConfig
	client config_client.IConfigClient
}

func New(cfg *NacosConfig) *NacosReader {

	//create ServerConfig
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(cfg.IP, cfg.Port),
	}

	//create ClientConfig
	cc := constant.NewClientConfig(
		constant.WithUsername(cfg.UserName),
		constant.WithPassword(cfg.Password),
		constant.WithNamespaceId(cfg.NameSpaceId),
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithCacheDir(cfg.RuntimeDir+"/cache"),
		constant.WithLogDir(cfg.RuntimeDir+"/logs"),
		constant.WithLogLevel(cfg.LogLevel),
	)

	// create config client
	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		panic(err)
	}

	return &NacosReader{
		cfg:    cfg,
		client: client,
	}
}

func (n *NacosReader) GetConfig() (string, error) {
	//get config
	content, err := n.client.GetConfig(vo.ConfigParam{
		DataId: n.cfg.DataId,
		Group:  n.cfg.Group,
	})
	if err != nil {
		return "", err
	}

	logger.GetLogger().Info("nacos get config:\n" + content)
	return content, nil
}

func (n *NacosReader) AddListener(listener func(content string) error) error {

	//Listen config change,key=dataId+group+namespaceId.
	err := n.client.ListenConfig(vo.ConfigParam{
		DataId: n.cfg.DataId,
		Group:  n.cfg.Group,
		OnChange: func(namespace, group, dataId, data string) {
			logger.GetLogger().Info("nacos config changed group:" + group + ", dataId:" + dataId + ", content:" + data)
			if err := listener(data); err != nil {
				logger.GetLogger().Error("nacos config changed reload failed")
			}
		},
	})
	if err != nil {
		return err
	}

	return nil
}
