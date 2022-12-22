package athena

import (
	"encoding/json"
	"fmt"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/spf13/viper"
	"golang.org/x/text/language"
	"log"
	"path/filepath"
	"strings"
	"sync"
)

// I18nOpt 国际化配置
type I18nOpt struct {
	DefaultLanguage string `mapstructure:"defaultLanguage"` // 默认语言（如：zh en ja）
}

// I18nConf 限流规则map
type I18nConf struct {
	I18nOpt *I18nOpt `mapstructure:"i18n"`
}

func (this *I18nConf) InitDefaultConfig(vp *viper.Viper) {}

var i18nModule *I18nModule

type I18nModule struct {
	conf     I18nConf
	Localize map[string]*i18n.Localizer
	Bundle   *i18n.Bundle
}

func NewI18nModule() *I18nModule {
	if i18nModule == nil {
		i18nModule = &I18nModule{Localize: make(map[string]*i18n.Localizer)}
	}
	return i18nModule
}

func (this *I18nModule) initConf() {
	once := sync.Once{}
	once.Do(func() {
		AddViperUnmarshal(&this.conf, nil)

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
	localize, exist := i18nModule.Localize[i18nModule.conf.I18nOpt.DefaultLanguage]
	if exist {
		return localize, nil
	}
	return localize, fmt.Errorf("invalid default localizer: %s", i18nModule.conf.I18nOpt.DefaultLanguage)
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
