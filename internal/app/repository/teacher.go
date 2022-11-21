package repository

import (
	"context"

	"onthemat/internal/app/model"
	"onthemat/pkg/ent"

	"onthemat/pkg/ent/teacher"

	"onthemat/pkg/entx"
)

type TeacherRepository interface {
	Create(ctx context.Context, t *ent.Teacher, userId int) error
	List(ctx context.Context, yogaSorts []*string, areas []*string) ([]*ent.Teacher, error)
}

type teacherRepository struct {
	db *ent.Client
}

func NewTeacherRepository(db *ent.Client) TeacherRepository {
	return &teacherRepository{
		db: db,
	}
}

func (repo *teacherRepository) Create(ctx context.Context, t *ent.Teacher, userId int) error {
	return entx.WithTx(ctx, repo.db, func(tx *ent.Tx) (err error) {
		if err = repo.db.Teacher.Create().
			SetAge(t.Age).
			SetName(t.Name).
			SetUserID(userId).
			Exec(ctx); err != nil {
			return
		}

		if err = repo.db.User.UpdateOneID(userId).SetType(model.AcademyType).Exec(ctx); err != nil {
			return
		}

		return
	})
}

func (repo *teacherRepository) Update(ctx context.Context, t *ent.Teacher, userId int) error {
	return repo.db.Teacher.Update().
		SetAge(t.Age).
		SetName(t.Name).
		SetIsProfileOpen(t.IsProfileOpen).
		Where(teacher.IDEQ(userId)).
		Exec(ctx)
}

func (repo *teacherRepository) Get(ctx context.Context, userId int) (*ent.Teacher, error) {
	return repo.db.Teacher.Query().
		Where(teacher.ID(userId)).
		Only(ctx)
}

func (repo *teacherRepository) List(ctx context.Context, yogaSorts []*string, areas []*string) ([]*ent.Teacher, error) {
	clause := repo.db.Debug().Teacher.Query()
	// clause = repo.conditionQuery(yogaSorts, areas, clause)
	return clause.All(ctx)
}
