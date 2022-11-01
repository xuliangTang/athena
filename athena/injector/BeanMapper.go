package injector

import (
	"reflect"
)

type BeanMapper map[reflect.Type]reflect.Value

func (this BeanMapper) add(bean any) {
	t := reflect.TypeOf(bean)
	if t.Kind() == reflect.Ptr {
		this[t] = reflect.ValueOf(bean)
	}
}

func (this BeanMapper) get(bean any) reflect.Value {
	var t reflect.Type
	if bt, ok := bean.(reflect.Type); ok {
		t = bt
	} else {
		t = reflect.TypeOf(bt)
	}

	if v, ok := this[t]; ok {
		return v
	}

	// 处理interface
	if t.Kind() == reflect.Interface {
		for k, v := range this {
			if k.Implements(t) {
				return v
			}
		}
	}

	return reflect.Value{}
}
