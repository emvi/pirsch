package pirsch

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEventOptions_getMetaData(t *testing.T) {
	options := EventOptions{
		Meta: map[string]string{
			"key":   "value",
			"hello": "world",
		},
	}
	k, v := options.getMetaData()
	assert.Len(t, k, 2)
	assert.Len(t, v, 2)
	assert.Contains(t, k, "key")
	assert.Contains(t, k, "hello")
	assert.Contains(t, v, "value")
	assert.Contains(t, v, "world")
}
