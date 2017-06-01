package tools

import (
	"math"
	"github.com/zalando/skipper/filters"
)

func ParseEskipFloatArg(arg interface{}) (float64, error) {
	if number, ok := arg.(float64); ok {
		return float64(number), nil
	} else {
		return 0, filters.ErrInvalidFilterParameters
	}
}

func ParseEskipIntArg(arg interface{}) (int, error) {
	if number, ok := arg.(float64); ok && math.Trunc(number) == number {
		return int(number), nil
	} else {
		return 0, filters.ErrInvalidFilterParameters
	}
}

func ParseEskipUint8Arg(arg interface{}) (uint8, error) {
	if number, ok := arg.(float64); ok && math.Trunc(number) == number {
		return uint8(number), nil
	} else {
		return 0, filters.ErrInvalidFilterParameters
	}
}

func ParseEskipStringArg(arg interface{}) (string, error) {
	if str, ok := arg.(string); ok {
		return string(str), nil
	} else {
		return "", filters.ErrInvalidFilterParameters
	}
}