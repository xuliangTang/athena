package config

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
	FileName   string
	AppPath    string
	Port       int
	Logging    *LoggingOpt
	Cors       *CorsOpt
	ErrorCache *ErrorCacheOpt
}

type CorsOpt struct {
	Enable bool
}

type ErrorCacheOpt struct {
	Enable bool
}

type LoggingOpt struct {
	RequestLogEnable bool // 是否开启http访问日志
	LogAccess        *LogFileOpt
	LogError         *LogFileOpt
}

type LogFileOpt struct {
	FilePath   string // 输出文件目录
	MaxSize    int    // 日志大小限制，单位MB
	MaxAge     int    // 历史日志文件保留天数
	MaxBackups int    // 最大保留历史日志数量
}

func (this *AppConfImpl) InitDefaultConfig(vp *viper.Viper) {
	vp.SetDefault("port", 80)
	vp.SetDefault("logging.requestLogEnable", false)
	vp.SetDefault("logging.logAccess.filePath", "/storage/logs/access.log")
	vp.SetDefault("logging.logAccess.maxSize", 255)
	vp.SetDefault("logging.logAccess.maxAge", 60)
	vp.SetDefault("logging.logAccess.maxBackups", 5)
	vp.SetDefault("logging.logError.filePath", "/storage/logs/error.log")
	vp.SetDefault("logging.logError.maxSize", 255)
	vp.SetDefault("logging.logError.maxAge", 180)
	vp.SetDefault("logging.logError.maxBackups", 5)
	vp.SetDefault("cors.enable", false)
	vp.SetDefault("errorCache.enable", false)
}

func init() {
	AppConf = AppConfImpl{
		AppPath:  Helper.GetWorkDir(),
		FileName: "application.yml",
	}

	AddViperUnmarshal(AppConf.FileName, &AppConf, func(vp *viper.Viper) OnConfigChangeRunFn {
		return func(in fsnotify.Event) {
			// 配置变更后重新解析
			if err := vp.Unmarshal(&AppConf); err != nil {
				log.Println(fmt.Sprintf("unmarshal conf failed: %s", err.Error()))
			}
		}
	})
}

type IConfig interface {
	// InitDefaultConfig 初始化配置默认值
	InitDefaultConfig(vp *viper.Viper)
}

type OnConfigChangeFn func(vp *viper.Viper) OnConfigChangeRunFn
type OnConfigChangeRunFn func(in fsnotify.Event)

func AddViperUnmarshal(fileName string, conf IConfig, onChange OnConfigChangeFn, decoderConfigOpts ...viper.DecoderConfigOption) {
	vp := viper.New()
	vp.SetConfigFile(fmt.Sprintf("%s/%s", AppConf.AppPath, fileName))
	vp.AutomaticEnv()
	if err := vp.ReadInConfig(); err != nil {
		log.Fatalln(fmt.Sprintf("error reading config file, %s", err))
	}

	// 初始化默认值
	conf.InitDefaultConfig(vp)

	if err := vp.Unmarshal(conf, decoderConfigOpts...); err != nil {
		log.Fatalln(fmt.Sprintf("unable to decode into struct, %v", err))
	}

	// 监控配置文件变化
	if onChange != nil {
		vp.WatchConfig()
		vp.OnConfigChange(onChange(vp))
	}
}
