package repository

import (
	"context"

	"onthemat/pkg/ent"
)

type TeacherRepository interface{}

type teacherRepository struct {
	db *ent.Client
}

func NewTeacherRepository(db *ent.Client) TeacherRepository {
	return &teacherRepository{
		db: db,
	}
}

func (repo *teacherRepository) Get(ctx context.Context) ([]*ent.TeacherWorkExperience, error) {
	return repo.db.TeacherWorkExperience.Query().All(ctx)
}
