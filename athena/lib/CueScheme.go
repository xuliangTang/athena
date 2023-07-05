package lib

import (
	"cuelang.org/go/cue"
	"fmt"
)

type CueScheme struct {
	schemes map[string]cue.Value
}

func NewCueScheme() *CueScheme {
	return &CueScheme{schemes: make(map[string]cue.Value)}
}

// AddScheme 注册scheme
func (this *CueScheme) AddScheme(name string, v cue.Value) {
	if v.Err() == nil {
		this.schemes[name] = v
	}
}

// GetScheme 获取cue.value
func (this *CueScheme) GetScheme(name string) (cue.Value, error) {
	if v, ok := this.schemes[name]; ok {
		return v, nil
	}
	return cue.Value{}, fmt.Errorf("not found scheme %s", name)
}

// MustGetScheme 获取cue.value
func (this *CueScheme) MustGetScheme(name string) cue.Value {
	if v, ok := this.schemes[name]; ok {
		return v
	}
	panic("Not Found Scheme " + name)
}
