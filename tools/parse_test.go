package tools

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestParseEskipIntArgSuccess(t *testing.T) {
	result, _ := ParseEskipIntArg(1.0)
	assert.Equal(t, 1, result)
}

func TestParseEskipIntArgFailure(t *testing.T) {
	_, err := ParseEskipIntArg(1.2)
	assert.NotNil(t, err, "There should be an error")
}

func TestParseEskipUint8ArgSuccess(t *testing.T) {
	result, _ := ParseEskipIntArg(1.0)
	assert.Equal(t, 1, result)
}

func TestParseEskipUint8ArgFailure(t *testing.T) {
	_, err := ParseEskipIntArg(1.2)
	assert.NotNil(t, err, "There should be an error")
}

func TestParseEskipStringArgSuccess(t *testing.T) {
	result, _ := ParseEskipStringArg("jpeg")
	assert.Equal(t, "jpeg", result)
}

func TestParseEskipStringArgFailure(t *testing.T) {
	_, err := ParseEskipStringArg(1.2)
	assert.NotNil(t, err, "There should be an error")
}

func TestParseEskipFloatArgSuccess(t *testing.T) {
	result, _ := ParseEskipFloatArg(54.321)
	assert.Equal(t, 54.321, result)
}

func TestParseEskipFloatArgFailure(t *testing.T) {
	_, err := ParseEskipFloatArg("1")
	assert.NotNil(t, err)
}