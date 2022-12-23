package configurations

import (
	"github.com/xuliangTang/athena/tests/core"
)

// 注入 回调handler
type K8sHandler struct{}

func NewK8sHandler() *K8sHandler {
	return &K8sHandler{}
}

// deployment handler
func (this *K8sHandler) DepHandlers() *core.DepHandler {
	return &core.DepHandler{}
}

//func(this *K8sHandler) PodHandlers() *core.DepHandler{
//	return &core.DepHandler{}
//}
