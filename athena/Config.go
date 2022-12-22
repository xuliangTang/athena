package athena

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	Helper "github.com/xuliangTang/athena/cmd/helper"
	"log"
)

// AppConf 配置中心
var AppConf AppConfImpl

type AppConfImpl struct {
	AppPath   string
	Port      int
	LogAccess *LogOpt
	LogError  *LogOpt
}

type LogOpt struct {
	FilePath   string // 输出文件目录
	MaxSize    int    // 日志大小限制，单位MB
	MaxAge     int    // 历史日志文件保留天数
	MaxBackups int    // 最大保留历史日志数量
}

func (this *AppConfImpl) InitDefaultConfig(vp *viper.Viper) {
	vp.SetDefault("port", 80)
	vp.SetDefault("logAccess.filePath", "/storage/logs/access.log")
	vp.SetDefault("logAccess.maxSize", 255)
	vp.SetDefault("logAccess.maxAge", 60)
	vp.SetDefault("logAccess.maxBackups", 5)
	vp.SetDefault("logError.filePath", "/storage/logs/error.log")
	vp.SetDefault("logError.maxSize", 255)
	vp.SetDefault("logError.maxAge", 180)
	vp.SetDefault("logError.maxBackups", 5)
}

func init() {
	AppConf = AppConfImpl{
		AppPath: Helper.GetWorkDir(),
	}

	AddViperUnmarshal(&AppConf, func(vp *viper.Viper) OnConfigChangeRunFn {
		return func(in fsnotify.Event) {
			// 配置变更后重新解析
			if err := vp.Unmarshal(&AppConf); err != nil {
				log.Println(fmt.Sprintf("unmarshal conf failed: %s", err.Error()))
			}
		}
	})
}

// ConfigModule 项目配置模块
type ConfigModule struct {
	ConfigInit IConfig
}

type IConfig interface {
	// InitDefaultConfig 初始化配置默认值
	InitDefaultConfig(vp *viper.Viper)
}

// NewConfigModule 自定义配置映射模块
func NewConfigModule(init IConfig) *ConfigModule {
	return &ConfigModule{
		ConfigInit: init,
	}
}

func (this *ConfigModule) Run() error {
	AddViperUnmarshal(this.ConfigInit, func(vp *viper.Viper) OnConfigChangeRunFn {
		return func(in fsnotify.Event) {
			// 配置变更后重新解析
			if err := vp.Unmarshal(&this.ConfigInit); err != nil {
				log.Println(fmt.Sprintf("unmarshal conf failed: %s", err.Error()))
			}
		}
	})

	/*viper.SetConfigName("application")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading config file, %s", err)
	}

	// 初始化默认值
	this.ConfigInit.InitDefaultConfig(viper.New())

	if err := viper.Unmarshal(this.ConfigInit); err != nil {
		return fmt.Errorf("unable to decode into struct, %v", err)
	}

	// 监控配置文件变化
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		// 配置变更后重新解析
		if err := viper.Unmarshal(this.ConfigInit); err != nil {
			log.Println(fmt.Sprintf("unmarshal conf failed: %s", err.Error()))
		}
	})*/

	return nil
}

type OnConfigChangeFn func(vp *viper.Viper) OnConfigChangeRunFn
type OnConfigChangeRunFn func(in fsnotify.Event)

func AddViperUnmarshal(conf IConfig, onChange OnConfigChangeFn) {
	vp := viper.New()
	vp.SetConfigFile(AppConf.AppPath + "/application.yml")
	vp.AutomaticEnv()
	if err := vp.ReadInConfig(); err != nil {
		log.Fatalln(fmt.Sprintf("error reading config file, %s", err))
	}

	// 初始化默认值
	conf.InitDefaultConfig(vp)

	if err := vp.Unmarshal(conf); err != nil {
		log.Fatalln(fmt.Sprintf("unable to decode into struct, %v", err))
	}

	// 监控配置文件变化
	if onChange != nil {
		vp.WatchConfig()
		vp.OnConfigChange(onChange(vp))
	}
}
