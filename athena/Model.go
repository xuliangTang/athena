package athena

import (
	"encoding/json"
	"log"
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
