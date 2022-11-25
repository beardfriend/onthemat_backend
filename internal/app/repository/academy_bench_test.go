package repository

import (
	"context"
	"testing"

	"onthemat/internal/app/config"
	"onthemat/internal/app/infrastructure"
	"onthemat/pkg/ent"
)

func BenchmarkAcademy(b *testing.B) {
	c := config.NewConfig()
	if err := c.Load("../../../configs"); err != nil {
		b.Fail()
		return
	}
	db := infrastructure.NewPostgresDB(c)
	aca := NewAcademyRepository(db)

	data := &ent.Academy{
		ID:            10,
		Name:          "name",
		CallNumber:    "010432266333",
		AddressRoad:   "상세주소",
		AddressDetail: nil,
		SigunguID:     20,
	}

	for i := 0; i < b.N; i++ {
		err := aca.Patch(context.Background(), data)
		b.Log(err)
	}
}
