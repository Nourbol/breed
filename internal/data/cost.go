package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Cost int32

var ErrInvalidRuntimeFormat = errors.New("invalid runtime format")

func (c Cost) MarshalJSON() ([]byte, error) {
	jsonValue := fmt.Sprintf("%d dollars", c)
	quotedJSONValue := strconv.Quote(jsonValue)
	return []byte(quotedJSONValue), nil
}

func (c *Cost) UnmarshalJSON(jsonValue []byte) error {
	unquotedJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	parts := strings.Split(unquotedJSONValue, " ")
	if len(parts) != 2 || parts[1] != "$" {
		return ErrInvalidRuntimeFormat
	}

	sign, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	*c = Cost(sign)
	return nil
}
