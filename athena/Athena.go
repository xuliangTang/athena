package athena

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xuliangTang/athena/athena/injector"
	"github.com/xuliangTang/athena/athena/task"
	"go.uber.org/zap"
	"log"
	"net/http"
)

type Athena struct {
	*gin.Engine
	g     *gin.RouterGroup
	props []any
}

func Ignite(fns ...FrameConfAttrFn) *Athena {
	FrameConfAttrFns(fns).apply(FrameConf)

	g := &Athena{Engine: gin.New()}
	g.Use(CorsHandler(), ErrorHandler(), RequestHandler())
	return g
}

// Launch 最终启动函数
func (this *Athena) Launch() {
	task.GetCron().Start()
	this.Run(fmt.Sprintf(":%d", FrameConf.Port))
}

func (this *Athena) Handle(httpMethod, relativePath string, handler interface{}) *Athena {
	if h := Convert(handler); h != nil {
		this.g.Handle(httpMethod, relativePath, h)
	}
	return this
}

// Load 初始化加载模块
func (this *Athena) Load(ls ...ILoad) *Athena {
	for _, l := range ls {
		err := l.Run()
		if err != nil {
			Logger().Error("load error",
				zap.String("info", err.Error()),
			)
			log.Fatalln(err.Error())
		}
	}
	return this
}

// Attach 加入全局中间件
func (this *Athena) Attach(f IFairing) *Athena {
	this.Use(func(context *gin.Context) {
		err := f.OnRequest(context)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			context.Next()
		}
	})
	return this
}

// Mount 挂载
func (this *Athena) Mount(group string, fs []IFairing, classes ...IClass) *Athena {
	if fs != nil && len(fs) > 0 {
		var handlers []gin.HandlerFunc
		for _, f := range fs {
			handlers = append(handlers, func(context *gin.Context) {
				err := f.OnRequest(context)
				if err != nil {
					context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				} else {
					context.Next()
				}
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

// Beans 依赖注入对象
func (this *Athena) Beans(beans ...any) *Athena {
	injector.BeanFactory.Set(beans...)
	this.props = append(this.props, beans...)
	return this
}

// CronTask 创建定时任务
func (this *Athena) CronTask(expr string, f func()) *Athena {
	_, err := task.GetCron().AddFunc(expr, f)
	if err != nil {
		Logger().Error("cron task error",
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
