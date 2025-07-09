package config

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
)

// 初始化本地配置并监听文件变化
func LoadConfigFromFile(configPath string, configType string) (*Config, error) {
	fmt.Println("Initializing local config from file:", configPath)

	// 初始化配置指针（后续更新会直接修改该指针指向的内存）
	cfg := &Config{}
	v := viper.New()

	v.SetConfigFile(configPath)
	v.SetConfigType(configType)
	// 修改映射字段为json tag
	hook := func(dc *mapstructure.DecoderConfig) {
		dc.TagName = "json"
	}

	// 首次读取本地文件
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read local config: %w", err)
	}

	// 首次解析到cfg指针
	if err := v.Unmarshal(cfg, hook); err != nil {
		return nil, fmt.Errorf("failed to unmarshal local config: %w", err)
	}

	// 监听本地文件变化：更新时直接修改cfg指向的内存
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("Local config file changed: %s\n", e.Name)
		// 重新解析到原指针（关键：不创建新指针，直接更新数据）
		if err := v.Unmarshal(cfg, hook); err != nil {
			fmt.Printf("failed to update local config: %v\n", err)
			return
		}
		fmt.Println("Local config updated successfully")
	})

	fmt.Println("Local config initialization completed")
	fmt.Printf("%+v\n", cfg)
	return cfg, nil // 返回指针，外部持有后会同步更新
}
