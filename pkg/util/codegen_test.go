package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandShortlyCode(t *testing.T) {
	assert.Len(t, RandShortlyCode(6), 6)
}
