package injector

import (
	"reflect"
)

var BeanFactory *BeanFactoryImpl

func init() {
	BeanFactory = NewBeanFactoryImpl()
}

type BeanFactoryImpl struct {
	beanMapper BeanMapper
}

func (this *BeanFactoryImpl) Set(beans ...any) {
	if beans == nil || len(beans) == 0 {
		return
	}

	for _, bean := range beans {
		this.beanMapper.add(bean)
	}
}

func (this *BeanFactoryImpl) Get(k any) any {
	if k == nil {
		return nil
	}

	getBean := this.beanMapper.get(k)
	if getBean.IsValid() {
		return getBean.Interface()
	}

	return nil
}

// Configuration 扫描配置类
func (this *BeanFactoryImpl) Configuration(cfgs ...any) {
	for _, cfg := range cfgs {
		t := reflect.TypeOf(cfg)
		if t.Kind() != reflect.Ptr {
			panic("configurations required ptr object")
		}
		if t.Elem().Kind() != reflect.Struct {
			continue
		}

		this.Set(cfg)
		this.Inject(cfg)

		v := reflect.ValueOf(cfg)
		for i := 0; i < t.NumMethod(); i++ {
			method := v.Method(i)
			callRet := method.Call(nil)

			if callRet != nil && len(callRet) == 1 {
				this.Set(callRet[0].Interface())
			}
		}
	}
}

// Inject 依赖注入bean
func (this *BeanFactoryImpl) Inject(cls any) {
	if cls == nil {
		return
	}

	v := reflect.ValueOf(cls)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		if v.Field(i).CanSet() && field.Tag.Get("inject") != "" {
			if field.Tag.Get("inject") == "-" { // 单例
				if getV := this.Get(field.Type); getV != nil {
					v.Field(i).Set(reflect.ValueOf(getV))
					this.Inject(getV) // 递归处理循环依赖
				}
			}
		}
	}
}

func NewBeanFactoryImpl() *BeanFactoryImpl {
	return &BeanFactoryImpl{beanMapper: make(BeanMapper)}
}

func (this *BeanFactoryImpl) GetBeanMapper() BeanMapper {
	return this.beanMapper
}
