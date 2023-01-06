package athena

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/xuliangTang/athena/athena/config"
	"github.com/xuliangTang/athena/athena/injector"
	"github.com/xuliangTang/athena/athena/interfaces"
	"github.com/xuliangTang/athena/athena/lib"
	"github.com/xuliangTang/athena/athena/middlewares"
	"github.com/xuliangTang/athena/athena/task"
	"go.uber.org/zap"
	"log"
	"reflect"
)

type Athena struct {
	*gin.Engine
	g     *gin.RouterGroup
	props []any
}

func Ignite() *Athena {
	g := &Athena{Engine: gin.New()}
	g.registerSysMiddleware()
	return g
}

// Launch 最终启动函数
func (this *Athena) Launch() {
	this.applyAll()
	task.GetCron().Start()
	this.Run(fmt.Sprintf(":%d", config.AppConf.Port))
}

func (this *Athena) Handle(httpMethod, relativePath string, handler interface{}) *Athena {
	if h := Convert(handler); h != nil {
		this.g.Handle(httpMethod, relativePath, h)
	}
	return this
}

// MappingConfig 映射配置文件到实体对象中
func (this *Athena) MappingConfig(entity config.IConfig) *Athena {
	config.AddViperUnmarshal(config.AppConf.FileName, entity, func(vp *viper.Viper) config.OnConfigChangeRunFn {
		return func(in fsnotify.Event) {
			// 配置变更后重新解析
			if err := vp.Unmarshal(&entity); err != nil {
				log.Println(fmt.Sprintf("unmarshal config failed: %s", err.Error()))
			}
		}
	})

	return this
}

// RegisterPlugin 注册插件
func (this *Athena) RegisterPlugin(plugins ...interfaces.IPlugin) *Athena {
	for _, plugin := range plugins {
		if !plugin.Enabler() {
			continue
		}

		plugin.InitModule()
	}

	return this
}

// 根据开关开启中间件
func (this *Athena) registerSysMiddleware() {
	if config.AppConf.Cors.Enable {
		this.Attach(middlewares.NewCors())
	}

	if config.AppConf.ErrorCache.Enable {
		this.Attach(middlewares.NewErrorCache())
	}

	if config.AppConf.Logging.RequestLogEnable {
		this.Attach(middlewares.NewRequestLog())
	}
}

// Attach 加入全局中间件
func (this *Athena) Attach(f interfaces.IFairing) *Athena {
	this.Use(func(context *gin.Context) {
		f.OnRequest(context)
	})
	return this
}

// Mount 挂载
func (this *Athena) Mount(group string, fs []interfaces.IFairing, classes ...IClass) *Athena {
	if fs != nil && len(fs) > 0 {
		var handlers []gin.HandlerFunc
		for _, f := range fs {
			handlers = append(handlers, func(context *gin.Context) {
				f.OnRequest(context)
			})
		}
		this.g = this.Group(group, handlers...)
	} else {
		this.g = this.Group(group)
	}

	for _, class := range classes {
		class.Build(this)
		// this.setProp(class)
		injector.BeanFactory.Inject(class)
	}
	return this
}

// Configuration 定义配置类，会被自动扫描注册到bean对象
func (this *Athena) Configuration(cfgs ...any) *Athena {
	injector.BeanFactory.Configuration(cfgs...)
	return this
}

// Beans 依赖注入对象
func (this *Athena) Beans(beans ...any) *Athena {
	injector.BeanFactory.Set(beans...)
	// this.props = append(this.props, beans...)
	return this
}

func (this *Athena) applyAll() {
	for t, v := range injector.BeanFactory.GetBeanMapper() {
		if t.Elem().Kind() == reflect.Struct {
			injector.BeanFactory.Inject(v.Interface())
		}
	}
}

// CronTask 创建定时任务
func (this *Athena) CronTask(expr string, f func()) *Athena {
	_, err := task.GetCron().AddFunc(expr, f)
	if err != nil {
		lib.Logger().Error("cron task error",
			zap.String("expr", expr),
			zap.String("info", err.Error()),
		)
	}
	return this
}

// 获取注入对象
/*func (this *Athena) getProp(t reflect.Type) any {
	for _, prop := range this.props {
		if t == reflect.TypeOf(prop) {
			return prop
		}
	}
	return nil
}*/

// 基于指针结构体属性的依赖注入
/*func (this *Athena) setProp(class IClass) {
	vClass := reflect.ValueOf(class).Elem()
	for i := 0; i < vClass.NumField(); i++ {
		field := vClass.Field(i)
		if !field.CanSet() || !field.IsNil() || field.Kind() != reflect.Ptr {
			continue
		}
		if prop := this.getProp(field.Type()); prop != nil {
			field.Set(reflect.New(field.Type().Elem()))
			field.Elem().Set(reflect.ValueOf(prop).Elem())
		}
	}
}*/
