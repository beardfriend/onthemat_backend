package repository

import (
	"context"

	"onthemat/internal/app/model"
	"onthemat/pkg/ent"
	"onthemat/pkg/ent/teacher"
	"onthemat/pkg/ent/user"

	"onthemat/pkg/entx"
)

type TeacherRepository interface {
	Create(ctx context.Context, d *ent.Teacher) error
	GetOnlyIdByUserId(ctx context.Context, userId int) (id int, err error)
}

type teacherRepository struct {
	db *ent.Client
}

func NewTeacherRepository(db *ent.Client) TeacherRepository {
	return &teacherRepository{
		db: db,
	}
}

func (repo *teacherRepository) Create(ctx context.Context, d *ent.Teacher) error {
	return entx.WithTx(ctx, repo.db, func(tx *ent.Tx) (err error) {
		clause := repo.db.Teacher.Create().
			SetUserID(d.UserID).
			SetNillableProfileImageURL(d.ProfileImageURL).
			SetName(d.Name).
			SetNillableAge(d.Age).
			SetNillableIntroduce(d.Introduce)

		if len(d.Edges.Yoga) > 0 {
			clause.AddYoga(d.Edges.Yoga...)
		}

		if len(d.Edges.Sigungu) > 0 {
			clause.AddSigungu(d.Edges.Sigungu...)
		}

		if len(d.Edges.WorkExperience) > 0 {
			clause.AddWorkExperience(d.Edges.WorkExperience...)
		}

		if len(d.Edges.YogaRaw) > 0 {
			clause.AddYogaRaw(d.Edges.YogaRaw...)
		}

		if err = clause.Exec(ctx); err != nil {
			return
		}

		err = tx.User.Update().
			SetType(model.TeacherType).
			Where(user.IDEQ(d.UserID)).
			Exec(ctx)
		if err != nil {
			return
		}

		return
	})
}

func (repo *teacherRepository) GetOnlyIdByUserId(ctx context.Context, userId int) (id int, err error) {
	return repo.db.Teacher.Query().Where(teacher.UserIDEQ(userId)).OnlyID(ctx)
}
