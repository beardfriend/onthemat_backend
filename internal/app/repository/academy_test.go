package repository

import (
	"context"
	"testing"

	"onthemat/internal/app/infrastructor"
)

func TestAcadmeyCreate(t *testing.T) {
	db := infrastructor.NewPostgresDB()
	repo := NewAcademyRepository(db)
	if err := repo.Create(context.Background()); err != nil {
		t.Error(err)
	}
}
