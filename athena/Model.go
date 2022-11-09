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

