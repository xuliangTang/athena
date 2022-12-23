package configurations

import (
	"github.com/xuliangTang/athena/tests/core"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
)

type K8sConfig struct {
	DepHandler *core.DepHandler `inject:"-"`
}

func NewK8sConfig() *K8sConfig {
	return &K8sConfig{}
}

func (*K8sConfig) InitClient() *kubernetes.Clientset {
	config := &rest.Config{
		Host:        "110.41.142.160:8009",
		BearerToken: "",
	}
	c, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	return c
}

// 初始化Informer
func (this *K8sConfig) InitInformer() informers.SharedInformerFactory {
	fact := informers.NewSharedInformerFactory(this.InitClient(), 0)

	depInformer := fact.Apps().V1().Deployments()
	depInformer.Informer().AddEventHandler(this.DepHandler)

	fact.Start(wait.NeverStop)

	return fact
}
