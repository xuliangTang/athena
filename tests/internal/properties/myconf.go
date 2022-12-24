package properties

import "github.com/spf13/viper"

type MyConfImpl struct {
	MyName string
	MyAge  int8
	Ex     *ExConfImpl
}

type ExConfImpl struct {
	ExName string
}

func (this *MyConfImpl) InitDefaultConfig(vp *viper.Viper) {
	vp.SetDefault("myAge", 99)
}

var MyConf MyConfImpl
