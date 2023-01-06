package plugins

import (
	"encoding/json"
	"fmt"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/spf13/viper"
	"github.com/xuliangTang/athena/athena/config"
	"golang.org/x/text/language"
	"log"
	"path/filepath"
	"strings"
	"sync"
)

// I18n @Plugin 国际化
type I18n struct {
	config   I18nConfig
	Localize map[string]*i18n.Localizer
	Bundle   *i18n.Bundle
}

func NewI18n() *I18n {
	if i18nModule == nil {
		i18nModule = &I18n{Localize: make(map[string]*i18n.Localizer)}
	}
	return i18nModule
}

var i18nModule *I18n

func (this *I18n) Enabler() bool {
	this.mappingConfig()
	return this.config.I18nOpt.Enable == true
}

func (this *I18n) mappingConfig() {
	once := sync.Once{}
	once.Do(func() {
		config.AddViperUnmarshal(config.AppConf.FileName, &this.config, nil)
	})
}

func (this *I18n) InitModule() {
	fs, err := filepath.Glob("lang/*")
	if err != nil {
		log.Fatalln("read lang dir error", err)
	}
	tag, err := language.Parse(this.config.I18nOpt.DefaultLanguage)
	if err != nil {
		log.Fatalln("parse language error: ", err)
	}
	this.Bundle = i18n.NewBundle(tag)
	this.Bundle.RegisterUnmarshalFunc("yaml", json.Unmarshal)
	replacer := strings.NewReplacer("lang/", "", "lang\\", "")
	for _, f := range fs {
		this.Bundle.MustLoadMessageFile(f)
		lang := replacer.Replace(strings.Split(f, ".")[0])
		this.Localize[lang] = i18n.NewLocalizer(this.Bundle, lang)
	}
}

// I18nOpt 国际化配置
type I18nOpt struct {
	Enable          bool
	DefaultLanguage string `mapstructure:"defaultLanguage"` // 默认语言（如：zh en ja）
}

// I18nConfig 限流规则map
type I18nConfig struct {
	I18nOpt I18nOpt `mapstructure:"i18n"`
}

func (this *I18nConfig) InitDefaultConfig(vp *viper.Viper) {
	vp.SetDefault("i18n.enable", true)
}

/*
func (this *I18nModule) initConf() {
	once := sync.Once{}
	once.Do(func() {
		athena.AddViperUnmarshal(&this.conf, nil)

		if this.conf.I18nOpt == nil {
			log.Fatalln("missing i18n conf")
		}
	})
}

func (this *I18nModule) Run() error {
	this.initConf()

	fs, err := filepath.Glob("lang/*")
	if err != nil {
		log.Fatalln("read lang dir error", err)
	}
	tag, err := language.Parse(this.conf.I18nOpt.DefaultLanguage)
	if err != nil {
		log.Fatalln("parse language error: ", err)
	}
	this.Bundle = i18n.NewBundle(tag)
	this.Bundle.RegisterUnmarshalFunc("yaml", json.Unmarshal)
	replacer := strings.NewReplacer("lang/", "", "lang\\", "")
	for _, f := range fs {
		this.Bundle.MustLoadMessageFile(f)
		lang := replacer.Replace(strings.Split(f, ".")[0])
		this.Localize[lang] = i18n.NewLocalizer(this.Bundle, lang)
	}

	return nil
}
*/

func checkIl8nModule() error {
	if i18nModule == nil {
		return fmt.Errorf("i18n module not loaded")
	}
	return nil
}

func GetDefaultLocalize() (*i18n.Localizer, error) {
	if err := checkIl8nModule(); err != nil {
		return nil, err
	}
	localize, exist := i18nModule.Localize[i18nModule.config.I18nOpt.DefaultLanguage]
	if exist {
		return localize, nil
	}
	return localize, fmt.Errorf("invalid default localizer: %s", i18nModule.config.I18nOpt.DefaultLanguage)
}

func GetLocalize(lang string) (*i18n.Localizer, error) {
	if err := checkIl8nModule(); err != nil {
		return nil, err
	}

	localize, exist := i18nModule.Localize[lang]
	if exist {
		return localize, nil
	}
	return localize, fmt.Errorf("invalid default localizer: %s", lang)
}
