package athena

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"
)

type Model interface {
	TableName() string
}
type Models string

func MakeModels(v interface{}) Models {
	b, err := json.Marshal(v)
	if err != nil {
		log.Println(err)
	}
	return Models(b)
}

// Conditions 自定义 where 条件
type Conditions struct {
	Query any
	Args  []any
}

func NewConditions(query any, args ...any) *Conditions {
	return &Conditions{Query: query, Args: args}
}

// NewConditionsWithQuery 根据注解生成 where 条件
func NewConditionsWithQuery(query any) *Conditions {
	var retQuery []string
	var retArgs []any

	v := reflect.ValueOf(query)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil
	}

	for i := 0; i < v.NumField(); i++ {
		typeField := v.Type().Field(i)
		field := v.Field(i)
		tagForm := typeField.Tag.Get("form")
		if tagForm != "" && !field.IsZero() {
			// 操作符默认=
			tagOp := typeField.Tag.Get("op")
			if tagOp == "" {
				tagOp = "="
			}

			// 字段名称优先使用column，没有的话使用form
			tagColumn := typeField.Tag.Get("column")
			if tagColumn != "" {
				tagForm = tagColumn
			}

			retQuery = append(retQuery, fmt.Sprintf("%s %s ?", tagForm, tagOp))

			switch tagOp {
			case "LIKE":
				retArgs = append(retArgs, fmt.Sprintf("%%%s%%", field.String()))
			default:
				retArgs = append(retArgs, field.Interface())
			}
		}
	}

	conditions := NewConditions(strings.Join(retQuery, " AND "), retArgs...)
	return conditions
}

// Preload 自定义预加载
type Preload struct {
	Query string
	Args  []any
}

func NewPreload(query string, args ...any) *Preload {
	return &Preload{Query: query, Args: args}
}

const DateTimeFormat = "2006-01-02 15:04:05"
const DateFormat = "2006-01-02"

// DateTime 自定义时间格式
type DateTime time.Time

func (t *DateTime) MarshalJSON() ([]byte, error) {
	tTime := time.Time(*t)
	return []byte(fmt.Sprintf("\"%v\"", tTime.Format(DateTimeFormat))), nil
}

func (t *DateTime) IsZero() bool {
	tm := time.Time(*t)
	return tm.IsZero()
}

func (t *DateTime) UnmarshalJSON(data []byte) (err error) {
	// 空值不进行解析
	if len(data) == 2 {
		*t = DateTime(time.Time{})
		return
	}

	// 指定解析的格式
	now, err := time.Parse(`"`+DateTimeFormat+`"`, string(data))
	*t = DateTime(now)
	return
}
