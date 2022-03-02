//go:build unit
// +build unit

package server

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServer_Log(t *testing.T) {
	l := ShortlyLog{
		Tag: "Tag",
	}

	l.Zap("Test")
	assert.Equal(t, l.format, "[Tag] Test")

	l.Zap("Test2")
	assert.Equal(t, l.format, "[Tag] Test2")

	l.Tag = "Warning"
	l.Zap("Warn")
	assert.Equal(t, l.format, "[Warning] Warn")

	l.ZapError("couldn't convert type", errors.New("invalid type"))
	assert.Equal(t, l.format, "[ERROR] couldn't convert type, invalid type")
}
