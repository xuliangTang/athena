package classes

import (
	"errors"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/xuliangTang/athena/athena"
	"github.com/xuliangTang/athena/athena/plugins"
	"github.com/xuliangTang/athena/tests/internal/properties"
	"golang.org/x/text/language"
	"net/http"
)

type TestClass struct {
}

func NewTestClass() *TestClass {
	return &TestClass{}
}

func (this *TestClass) test(ctx *gin.Context) athena.Json {
	return athena.Json{
		"message": "test",
		"my_name": properties.MyConf.MyName,
		"my_age":  properties.MyConf.MyAge,
		"ex_name": properties.MyConf.Ex.ExName,
	}
}

func (this *TestClass) ping(ctx *gin.Context) athena.Json {
	msg := "success"
	hystrix.Do("test1", func() error {
		resp, err := http.Get("https://www.google.com/")
		if err != nil || resp.StatusCode != http.StatusOK {
			fmt.Printf("请求失败:%v", err)
			return errors.New(fmt.Sprintf("error resp"))
		}
		return nil
	}, func(err error) error {
		if err != nil {
			fmt.Printf("circuitBreaker and err is %s\n", err.Error())
			msg = err.Error()
		}
		return nil
	})

	return athena.Json{"msg": msg}
}

func (this *TestClass) lang(ctx *gin.Context) athena.Json {
	localize := athena.Unwrap(plugins.GetDefaultLocalize()).(*i18n.Localizer)
	strDefault := athena.Unwrap(localize.Localize(&i18n.LocalizeConfig{
		MessageID:    "test.hello",
		TemplateData: map[string]string{"name": "Nick"},
	}))

	localizeEn := athena.Unwrap(plugins.GetLocalize(language.English.String())).(*i18n.Localizer)
	strEnglish := athena.Unwrap(localizeEn.Localize(&i18n.LocalizeConfig{
		MessageID:    "test.hello",
		TemplateData: map[string]string{"name": "Nick"},
	}))

	return athena.Json{"default": strDefault, "en": strEnglish}
}

func (this *TestClass) Build(athena *athena.Athena) {
	athena.Handle("GET", "/test", this.test)
	athena.Handle("GET", "/ping", this.ping)
	athena.Handle("GET", "/lang", this.lang)
}