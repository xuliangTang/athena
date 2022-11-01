package athena

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
)

// FrameConf 脚手架配置
var FrameConf *ConfImpl

type ConfImpl struct {
	AppPath string // 项目根目录
	Port    int
}

// IConfigInit 设置默认配置项的方法
type IConfigInit interface {
	InitDefaultConfig()
}

// ConfigModule 项目配置模块
type ConfigModule struct {
	AppConf    any
	ConfigInit IConfigInit
}

// NewConfigModule 在 load 方法中加载项目配置
func NewConfigModule(appConf any, init IConfigInit) *ConfigModule {
	return &ConfigModule{AppConf: appConf, ConfigInit: init}
}

func (this *ConfigModule) Run() error {
	viper.SetConfigName("application")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading config file, %s", err)
	}

	if err := viper.Unmarshal(&this.AppConf); err != nil {
		return fmt.Errorf("unable to decode into struct, %v", err)
	}

	if this.ConfigInit != nil {
		this.ConfigInit.InitDefaultConfig()
	}

	// 监控配置文件变化
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		if err := viper.Unmarshal(&this.AppConf); err != nil {
			log.Fatalln(fmt.Sprintf("unmarshal conf failed: %s", err.Error()))
		}
	})

	return nil
}
