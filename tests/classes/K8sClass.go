package classes

import (
	"github.com/gin-gonic/gin"
	"github.com/xuliangTang/athena/athena"
	"github.com/xuliangTang/athena/tests/core"
	"k8s.io/client-go/kubernetes"
)

type K8sClass struct {
	Client        *kubernetes.Clientset `inject:"-"`
	DeploymentMap *core.DeploymentMap   `inject:"-"`
}

func NewK8sClass() *K8sClass {
	return &K8sClass{}
}

func (this *K8sClass) deployments(ctx *gin.Context) *athena.Json {
	//deps, err := this.Client.AppsV1().Deployments("default").List(ctx, v1.ListOptions{})
	depList, err := this.DeploymentMap.ListByNS("default")
	athena.Error(err)
	return &athena.Json{"data": depList}
}

func (this *K8sClass) Build(athena *athena.Athena) {
	athena.Handle("GET", "deployments", this.deployments)
}
