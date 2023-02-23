package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Cost int32

var ErrInvalidCostFormat = errors.New("invalid cost format")

func (c Cost) MarshalJSON() ([]byte, error) {
	jsonValue := fmt.Sprintf("%d dollars", c)
	quotedJSONValue := strconv.Quote(jsonValue)
	return []byte(quotedJSONValue), nil
}

func (c *Cost) UnmarshalJSON(jsonValue []byte) error {
	unquotedJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidCostFormat
	}

	parts := strings.Split(unquotedJSONValue, " ")
	if len(parts) != 2 || parts[1] != "dollars" {
		return ErrInvalidCostFormat
	}

	sign, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return ErrInvalidCostFormat
	}

	*c = Cost(sign)
	return nil
}
