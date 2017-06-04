package parse

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestParseEskipIntArgSuccess(t *testing.T) {
	result, _ := EskipIntArg(1.0)
	assert.Equal(t, 1, result)
}

func TestParseEskipIntArgFailure(t *testing.T) {
	_, err := EskipIntArg(1.2)
	assert.NotNil(t, err, "There should be an error")
}

func TestParseEskipUint8ArgSuccess(t *testing.T) {
	result, _ := EskipIntArg(1.0)
	assert.Equal(t, 1, result)
}

func TestParseEskipUint8ArgFailure(t *testing.T) {
	_, err := EskipIntArg(1.2)
	assert.NotNil(t, err, "There should be an error")
}

func TestParseEskipStringArgSuccess(t *testing.T) {
	result, _ := EskipStringArg("jpeg")
	assert.Equal(t, "jpeg", result)
}

func TestParseEskipStringArgFailure(t *testing.T) {
	_, err := EskipStringArg(1.2)
	assert.NotNil(t, err, "There should be an error")
}

func TestParseEskipFloatArgSuccess(t *testing.T) {
	result, _ := EskipFloatArg(54.321)
	assert.Equal(t, 54.321, result)
}

func TestParseEskipFloatArgFailure(t *testing.T) {
	_, err := EskipFloatArg("1")
	assert.NotNil(t, err)
}

func TestEskipBoolArg(t *testing.T) {
	result, _ := EskipBoolArg(true)
	assert.True(t, result)
}

func TestEskipBoolArgFailure(t *testing.T) {
	_, err := EskipBoolArg(13)
	assert.NotNil(t, err)
}
