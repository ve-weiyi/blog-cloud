package cmd

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"

	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/config"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
	"github.com/ve-weiyi/ve-blog-golang/pkg/kit/logz"
)

// API命令参数
type ApiFlags struct {
	ConfigMode string // 运行方式 file|nacos
	ConfigType string // 配置类型（yaml/json等）

	// 本地配置相关
	LocalPath string // 本地配置文件路径（如config.yaml）

	// Nacos相关配置
	NacosHost       string // Nacos服务地址
	NacosPort       uint64 // Nacos服务端口
	NacosNamespace  string // Nacos命名空间
	NacosDataID     string // 配置DataID
	NacosGroup      string // 配置Group
	NacosRuntimeDir string // Nacos运行时目录
	NacosUsername   string // Nacos用户名
	NacosPassword   string // Nacos密码
}

var apiFlags = &ApiFlags{
	ConfigMode:      "file",
	ConfigType:      "yaml",
	LocalPath:       "config.yaml",
	NacosHost:       "veweiyi.cn",
	NacosPort:       8848,
	NacosNamespace:  "dev",
	NacosDataID:     "ve-blog-golang",
	NacosGroup:      "blog",
	NacosRuntimeDir: "runtime/nacos",
	NacosUsername:   "nacos",
	NacosPassword:   "nacos",
}

func NewApiCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "api",
		Short: "启动接口服务",
		Long:  `启动接口服务`,
		RunE:  runApi,
	}

	// 配置文件相关
	cmd.Flags().StringVarP(&apiFlags.ConfigMode, "config", "c", apiFlags.ConfigMode, "the way of read config file (file|nacos)")
	cmd.Flags().StringVarP(&apiFlags.LocalPath, "filepath", "f", apiFlags.LocalPath, "config file path")

	// Nacos 相关
	cmd.Flags().StringVar(&apiFlags.NacosHost, "n-host", apiFlags.NacosHost, "the host for nacos")
	cmd.Flags().Uint64Var(&apiFlags.NacosPort, "n-port", apiFlags.NacosPort, "the port for nacos")
	cmd.Flags().StringVar(&apiFlags.NacosNamespace, "n-namespace", apiFlags.NacosNamespace, "the namespace for nacos")
	cmd.Flags().StringVar(&apiFlags.NacosDataID, "n-data-id", apiFlags.NacosDataID, "the DataId for nacos")
	cmd.Flags().StringVar(&apiFlags.NacosGroup, "n-group", apiFlags.NacosGroup, "the group for nacos")
	cmd.Flags().StringVar(&apiFlags.NacosUsername, "n-user", apiFlags.NacosUsername, "the user for nacos")
	cmd.Flags().StringVar(&apiFlags.NacosPassword, "n-password", apiFlags.NacosPassword, "the password for nacos")

	return cmd
}

func runApi(cmd *cobra.Command, args []string) error {
	var err error
	var c *config.Config
	switch apiFlags.ConfigMode {
	case "file":
		c, err = config.LoadConfigFromFile(apiFlags.LocalPath, apiFlags.ConfigType)
	case "nacos":
		c, err = config.LoadConfigFromNacos(&config.LoadNacosConfigOption{
			ConfigType:      apiFlags.ConfigType,
			NacosHost:       apiFlags.NacosHost,
			NacosPort:       apiFlags.NacosPort,
			NacosNamespace:  apiFlags.NacosNamespace,
			NacosDataID:     apiFlags.NacosDataID,
			NacosGroup:      apiFlags.NacosGroup,
			NacosRuntimeDir: apiFlags.NacosRuntimeDir,
			NacosUsername:   apiFlags.NacosUsername,
			NacosPassword:   apiFlags.NacosPassword,
		})
	default:
		log.Fatalf("unsupported config file mode: %s", apiFlags.ConfigMode)
	}
	if err != nil {
		log.Fatalf("failed to initialize config: %v\n", err)
	}

	RunHttpServer(c)
	return nil
}

func RunHttpServer(c *config.Config) {
	// 初始化zap日志库
	logz.SetLog(&logz.LogConfig{
		Level:      c.Zap.Level,
		Mode:       c.Zap.Mode,
		Filename:   c.Zap.Filename,
		MaxSize:    c.Zap.MaxSize,
		MaxBackups: c.Zap.MaxBackups,
		MaxAge:     c.Zap.MaxAge,
		Compress:   c.Zap.Compress,
	})
	logz.L().Sugar().Infof("zap log init success. mode:%v, level:%v", c.Zap.Mode, c.Zap.Level)

	ctx := svctx.NewServiceContext(c)

	// 设置ReleaseMode则不会打印路由注册日志
	gin.SetMode(gin.DebugMode)
	engine := gin.Default()
	api.RegisterRouters(engine, ctx)

	logz.L().Sugar().Infof("register router success")

	address := fmt.Sprintf(":%d", c.System.Port)
	s := &http.Server{
		Addr:           address,
		Handler:        engine,
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	// 保证文本顺序输出
	// In order to ensure that the text order output can be deleted
	time.Sleep(10 * time.Microsecond)
	logz.L().Sugar().Infof("run server on http://localhost:%v success", c.System.Port)

	fmt.Printf(`
	欢迎使用 ve-blog-golang
	当前版本: %s
	微信号：wy791422171 QQ：791422171
	默认接口文档地址:http://localhost:%v/api/v1/swagger/index.html
`, c.System.Version, c.System.Port)
	fmt.Println(s.ListenAndServe().Error())
}
