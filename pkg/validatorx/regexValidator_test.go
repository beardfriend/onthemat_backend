package validatorx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImageContentTypeValidator(t *testing.T) {
	res1, _ := ImageExtensionValidator(".pngg")
	assert.Equal(t, res1, false)

	res2, _ := ImageExtensionValidator(".aa")
	assert.Equal(t, res2, false)

	res3, _ := ImageExtensionValidator(".png")
	assert.Equal(t, res3, true)
}
