package parse

import (
	"github.com/zalando/skipper/filters"
	"math"
)

// EskipFloatArg parse an eskip argument into a Float
func EskipFloatArg(arg interface{}) (float64, error) {
	if number, ok := arg.(float64); ok {
		return float64(number), nil
	}
	return 0, filters.ErrInvalidFilterParameters
}

// EskipIntArg parse an eskip argument into an Int
func EskipIntArg(arg interface{}) (int, error) {
	if number, ok := arg.(float64); ok && math.Trunc(number) == number {
		return int(number), nil
	}
	return 0, filters.ErrInvalidFilterParameters
}

// EskipUint8Arg parse an eskip argument into an UInt8
func EskipUint8Arg(arg interface{}) (uint8, error) {
	if number, ok := arg.(float64); ok && math.Trunc(number) == number {
		return uint8(number), nil
	}
	return 0, filters.ErrInvalidFilterParameters
}

// EskipStringArg parse an eskip argument into a String
func EskipStringArg(arg interface{}) (string, error) {
	if str, ok := arg.(string); ok {
		return string(str), nil
	}
	return "", filters.ErrInvalidFilterParameters
}

// EskipBoolArg parse an eskip argument into a Boolean
func EskipBoolArg(arg interface{}) (bool, error) {
	if value, ok := arg.(bool); ok {
		return value, nil
	}
	return false, filters.ErrInvalidFilterParameters
}
