package infrastructor

import (
	"testing"

	"onthemat/pkg/test"
)

func TestInitMariaDB(t *testing.T) {
	test.BeforeStart("../../../configs")

	db := NewMariaDB()

	if err := db.Raw(`SELECT 1`).Error; err != nil {
		t.Error(err)
	}
}
