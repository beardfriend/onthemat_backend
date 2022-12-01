package validatorx_test

import (
	"testing"

	"onthemat/pkg/validatorx"

	"github.com/stretchr/testify/assert"
)

func TestImageContentTypeValidator(t *testing.T) {
	res1, _ := validatorx.ImageExtensionValidator(".pngg")
	assert.Equal(t, res1, false)

	res2, _ := validatorx.ImageExtensionValidator(".aa")
	assert.Equal(t, res2, false)

	res3, _ := validatorx.ImageExtensionValidator(".png")
	assert.Equal(t, res3, true)
}
