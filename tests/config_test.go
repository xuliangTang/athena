package tests

import (
	"encoding/json"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/xuliangTang/athena/athena/config"
	"reflect"
	"testing"
)

type nodeConf struct {
	Nodes map[string]*nodeOpt `mapstructure:"nodes"`
}

type nodeOpt struct {
	Username string
	Password string
	Host     string
	Port     int
}

func (this *nodeConf) InitDefaultConfig(vp *viper.Viper) {
	vp.BindEnv("nodes", "NODES")
}

func (*nodeConf) jsonToNodesOptHookFunc() mapstructure.DecodeHookFuncType {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		// Check if the data type matches the expected one
		if f.Kind() != reflect.String {
			return data, nil
		}

		// Check that the target type is our custom type
		if t != reflect.TypeOf(map[string]*nodeOpt{}) {
			return data, nil
		}

		// Format/decode/parse the data and return the new value
		var m map[string]*nodeOpt
		json.Unmarshal([]byte(data.(string)), &m)
		return m, nil
	}
}

var node nodeConf

func Test_AddViperUnmarshal(t *testing.T) {
	config.AddViperUnmarshal("app.yml",
		&node,
		nil,
		viper.DecodeHook(node.jsonToNodesOptHookFunc()),
	)
	is := assert.New(t)
	is.NotNil(node.Nodes["node1"])
	is.NotNil(node.Nodes["node2"])
}
