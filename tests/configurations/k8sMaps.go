package configurations

import (
	"github.com/xuliangTang/athena/tests/core"
	"k8s.io/client-go/kubernetes"
)

type K8sMaps struct {
	K8sClient *kubernetes.Clientset `inject:"-"`
}

func NewK8sMaps() *K8sMaps {
	return &K8sMaps{}
}

// 初始化 deploymentmap
func (this *K8sMaps) InitDepMap() *core.DeploymentMap {
	return &core.DeploymentMap{}
}
