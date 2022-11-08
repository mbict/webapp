package router

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ExtractParamNames(t *testing.T) {
	assert.Nil(t, extractParamNames("/foo/bar/baz"))
	assert.Equal(t, []string{"bar"}, extractParamNames("/foo/{bar}/baz"))
	assert.Equal(t, []string{"bar", "baz"}, extractParamNames("/foo/{bar}/{baz}"))
	assert.Equal(t, []string{"bar", "wildcard"}, extractParamNames("/foo/{bar}/*wildcard"))
	assert.Equal(t, []string{"bar", "wildcard"}, extractParamNames("/foo/test@{bar}/*wildcard"))
}
