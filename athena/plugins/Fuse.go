package plugins

import (
	"github.com/afex/hystrix-go/hystrix"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"github.com/xuliangTang/athena/athena/config"
	"sync"
)

// Fuse @Plugin 熔断
type Fuse struct {
	config FuseConfig
}

func NewFuse() *Fuse {
	return &Fuse{}
}

func (this *Fuse) Enabler() bool {
	this.mappingConfig()
	return this.config.FuseOpt.Enable == true
}

func (this *Fuse) mappingConfig() {
	once := sync.Once{}
	once.Do(func() {
		config.AddViperUnmarshal(&this.config, func(vp *viper.Viper) config.OnConfigChangeRunFn {
			return func(in fsnotify.Event) {
				this.InitModule()
			}
		})
	})
}

func (this *Fuse) InitModule() {
	for k, v := range this.config.FuseOpt.Rules {
		this.addRuleByName(k, v)
	}
}

func (this *Fuse) addRuleByName(name string, opt *FuseRuleOpt) {
	hystrix.ConfigureCommand(name, hystrix.CommandConfig{
		Timeout:                opt.Timeout,
		MaxConcurrentRequests:  opt.MaxConcurrentRequests,
		SleepWindow:            opt.SleepWindow,
		RequestVolumeThreshold: opt.RequestVolumeThreshold,
		ErrorPercentThreshold:  opt.ErrorPercentThreshold,
	})
}

// FuseRuleOpt 熔断器配置
type FuseRuleOpt struct {
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

type FuseOpt struct {
	Enable bool
	Rules  map[string]*FuseRuleOpt `mapstructure:"fuseRules"`
}

type FuseConfig struct {
	FuseOpt FuseOpt `mapstructure:"fuse"`
}

func (this *FuseConfig) InitDefaultConfig(vp *viper.Viper) {
	vp.SetDefault("fuse.enable", true)
}

/*func (this *Fuse) getConf() *FuseRules {
	once := sync.Once{}
	rule := &FuseRules{}
	once.Do(func() {
		athena.AddViperUnmarshal(rule, func(vp *viper.Viper) athena.OnConfigChangeRunFn {
			return func(in fsnotify.Event) {
				this.addRuleByConf()
			}
		})
	})
	return rule
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
}*/
