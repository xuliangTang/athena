package athena

import (
	"encoding/json"
	"fmt"
	"log"
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

// DateTime 自定义时间格式
type DateTime time.Time

func (t *DateTime) MarshalJSON() ([]byte, error) {
	tTime := time.Time(*t)
	return []byte(fmt.Sprintf("\"%v\"", tTime.Format("2006-01-02 15:04:05"))), nil
}
