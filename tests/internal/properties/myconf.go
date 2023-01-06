package properties

import (
	"encoding/json"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"reflect"
)

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

var Nodes NodeConf

type NodeConf struct {
	Nodes map[string]*NodeOpt `mapstructure:"nodes"`
}

type NodeOpt struct {
	Username string
	Password string
	Host     string
	Port     int
}

func (this *NodeConf) InitDefaultConfig(vp *viper.Viper) {
	vp.BindEnv("nodes", "K8S_NODES")
}

// JsonToNodesOptHookFunc 解码json字符串的nodesOpt map
func (*NodeConf) JsonToNodesOptHookFunc() mapstructure.DecodeHookFuncType {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		// Check if the data type matches the expected one
		if f.Kind() != reflect.String {
			return data, nil
		}

		// Check that the target type is our custom type
		if t != reflect.TypeOf(map[string]*NodeOpt{}) {
			return data, nil
		}

		// Format/decode/parse the data and return the new value
		var m map[string]*NodeOpt
		json.Unmarshal([]byte(data.(string)), &m)
		return m, nil
	}
}
