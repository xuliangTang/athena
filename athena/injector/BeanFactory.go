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
