package infrastructor

import (
	"testing"

	"onthemat/internal/app/config"
)

func TestInitMariaDB(t *testing.T) {
	c := config.NewConfig()
	if err := c.Load("../../../configs"); err != nil {
		t.Error(err)
	}

	db := NewMariaDB(c)

	if err := db.Raw(`SELECT 1`).Error; err != nil {
		t.Error(err)
	}
}
