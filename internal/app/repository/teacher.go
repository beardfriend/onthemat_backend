package repository

import (
	"context"
	"database/sql/driver"

	"onthemat/internal/app/model"
	"onthemat/pkg/ent"

	"onthemat/pkg/ent/teacher"
	"onthemat/pkg/ent/teacherworkarea"
	"onthemat/pkg/ent/useryoga"
	"onthemat/pkg/entx"

	"entgo.io/ent/dialect/sql"
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
	clause = repo.conditionQuery(yogaSorts, areas, clause)
	return clause.All(ctx)
}

func (repo *teacherRepository) conditionQuery(yogaSorts []*string, areas []*string, clause *ent.TeacherQuery) *ent.TeacherQuery {
	if yogaSorts != nil {
		a := make([]driver.Value, len(yogaSorts))
		for i, v := range yogaSorts {
			a[i] = *v
		}
		clause = clause.Where(func(s *sql.Selector) {
			t := sql.Table(useryoga.Table)
			s.Join(t).On(s.C(teacher.FieldID), t.C(useryoga.FieldUserID))
			s.Where(sql.InValues(t.C(useryoga.FieldName), a...))
		})
	}
	if areas != nil {
		a := make([]driver.Value, len(areas))
		for i, v := range areas {
			a[i] = *v
		}
		clause = clause.Where(func(s *sql.Selector) {
			t := sql.Table(teacherworkarea.Table)
			s.Join(t).On(s.C(teacher.FieldID), t.C(teacherworkarea.FieldTeacherID))
			s.Where(sql.InValues(t.C(teacherworkarea.FieldGu), a...))
		})

	}

	return clause
}
