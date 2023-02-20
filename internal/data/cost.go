package data

import (
	"fmt"
	"strconv"
)

type Cost int32

func (r Cost) MarshalJSON() ([]byte, error) {
	jsonValue := fmt.Sprintf("%d dollars", r)
	quotedJSONValue := strconv.Quote(jsonValue)
	return []byte(quotedJSONValue), nil
}
