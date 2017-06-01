package parse

import (
	"math"
	"github.com/zalando/skipper/filters"
)

func EskipFloatArg(arg interface{}) (float64, error) {
	if number, ok := arg.(float64); ok {
		return float64(number), nil
	} else {
		return 0, filters.ErrInvalidFilterParameters
	}
}

func EskipIntArg(arg interface{}) (int, error) {
	if number, ok := arg.(float64); ok && math.Trunc(number) == number {
		return int(number), nil
	} else {
		return 0, filters.ErrInvalidFilterParameters
	}
}

func EskipUint8Arg(arg interface{}) (uint8, error) {
	if number, ok := arg.(float64); ok && math.Trunc(number) == number {
		return uint8(number), nil
	} else {
		return 0, filters.ErrInvalidFilterParameters
	}
}

func EskipStringArg(arg interface{}) (string, error) {
	if str, ok := arg.(string); ok {
		return string(str), nil
	} else {
		return "", filters.ErrInvalidFilterParameters
	}
}