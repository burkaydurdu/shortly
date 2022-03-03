package shortly

import (
	"testing"

	"github.com/burkaydurdu/shortly/internal/db"
	"github.com/stretchr/testify/assert"
)

func TestShortlyService_GenerateShortlyCode(t *testing.T) {
	var shortlyList = []db.Shortly{
		{
			Code: "XXXXXX",
		},
	}

	assert.Len(t, generateShortlyCode(shortlyList), 6)
}
