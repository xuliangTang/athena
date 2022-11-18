package athena

import (
	"github.com/afex/hystrix-go/hystrix"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"sync"
)

// FuseOpt 熔断器配置
type FuseOpt struct {
	// 定义执行command的超时时间(ms)
	Timeout int `mapstructure:"timeout"`
	// 定义command的最大并发量
	MaxConcurrentRequests int `mapstructure:"maxConcurrentRequests"`
	// 在熔断器被打开后，根据SleepWindow设置的时间控制多久后尝试服务是否可用
	SleepWindow int `mapstructure:"sleepWindow"`
	// 判断熔断开关的条件之一，统计10s内请求数量，达到这个请求数量后再根据错误率判断是否要开启熔断
	RequestVolumeThreshold int `mapstructure:"requestVolumeThreshold"`
	// 判断熔断开关的条件之一，错误率到达这个百分比后就会启动熔断
	ErrorPercentThreshold int `mapstructure:"errorPercentThreshold"`
}

type FuseRules struct {
	Rules map[string]*FuseOpt `mapstructure:"fuseRules"`
}

type Fuse struct {
}

func NewFuse() *Fuse {
	return &Fuse{}
}

func (this *Fuse) getConf() *FuseRules {
	once := sync.Once{}
	rule := &FuseRules{}
	once.Do(func() {
		vp := viper.New()
		vp.SetConfigFile(FrameConf.AppPath + "/application.yml")
		err := vp.ReadInConfig()
		if err != nil {
			log.Fatalln("read config.yaml error :", err)
		}
		errRule := vp.Unmarshal(&rule)
		if errRule != nil {
			log.Fatalln("unmarshal err :", errRule)
		}

		// 监控配置文件变化
		vp.WatchConfig()
		vp.OnConfigChange(func(in fsnotify.Event) {
			log.Println("更新了")
			this.addRuleByConf()
		})
	})
	return rule
}

func (this *Fuse) addRuleByName(name string, opt *FuseOpt) {
	hystrix.ConfigureCommand(name, hystrix.CommandConfig{
		Timeout:                opt.Timeout,
		MaxConcurrentRequests:  opt.MaxConcurrentRequests,
		SleepWindow:            opt.SleepWindow,
		RequestVolumeThreshold: opt.RequestVolumeThreshold,
		ErrorPercentThreshold:  opt.ErrorPercentThreshold,
	})
}

func (this *Fuse) addRuleByConf() {
	conf := this.getConf()
	for k, v := range conf.Rules {
		this.addRuleByName(k, v)
	}
}

func (this *Fuse) Run() error {
	this.addRuleByConf()
	return nil
}
