package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsURL(t *testing.T) {
	assert.Error(t, IsURL("asd"))
	assert.Error(t, IsURL("http"))
	assert.NoError(t, IsURL("http://"))
	assert.NoError(t, IsURL("http://sad"))
	assert.NoError(t, IsURL("https://burkaydurdu.github.io"))
	assert.Error(t, IsURL(""))
	assert.Error(t, IsURL("   "))
}
