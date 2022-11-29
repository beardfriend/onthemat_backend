package repository

import (
	"context"

	"onthemat/pkg/ent"
	"onthemat/pkg/ent/teacherworkexperience"
)

type TeacherWorkExperience interface {
	Create(ctx context.Context, value *ent.TeacherWorkExperience, teacherId int) error
	Update(ctx context.Context, value *ent.TeacherWorkExperience) error
	Delete(ctx context.Context, workExperienceId int) error
	Get(ctx context.Context, workExperienceId int) (*ent.TeacherWorkExperience, error)
	ListByTeacherID(ctx context.Context, teacherId int) ([]*ent.TeacherWorkExperience, error)
	CreateMany(ctx context.Context, value []*ent.TeacherWorkExperience, teacherId int) error
}

type teacherWorkExperience struct {
	db *ent.Client
}

func NewTeacherWorkExperience(db *ent.Client) TeacherWorkExperience {
	return &teacherWorkExperience{
		db: db,
	}
}

func (repo *teacherWorkExperience) Create(ctx context.Context, value *ent.TeacherWorkExperience, teacherId int) error {
	return repo.db.TeacherWorkExperience.Create().
		SetAcademyName(value.AcademyName).
		SetWorkStartAt(value.WorkStartAt).
		SetNillableWorkEndAt(value.WorkEndAt).
		SetNillableDescription(value.Description).
		SetTeacherID(teacherId).
		Exec(ctx)
}

func (repo *teacherWorkExperience) Get(ctx context.Context, workExperienceId int) (*ent.TeacherWorkExperience, error) {
	return repo.db.TeacherWorkExperience.Get(ctx, workExperienceId)
}

func (repo *teacherWorkExperience) ListByTeacherID(ctx context.Context, teacherId int) ([]*ent.TeacherWorkExperience, error) {
	return repo.db.Debug().TeacherWorkExperience.Query().Where(teacherworkexperience.TeacherIDEQ(teacherId)).
		All(ctx)
}

func (repo *teacherWorkExperience) Update(ctx context.Context, value *ent.TeacherWorkExperience) error {
	return repo.db.TeacherWorkExperience.
		Update().
		SetAcademyName(value.AcademyName).
		SetWorkStartAt(value.WorkStartAt).
		SetNillableWorkEndAt(value.WorkEndAt).
		SetNillableDescription(value.Description).
		SetTeacherID(value.ID).
		Exec(ctx)
}

func (repo *teacherWorkExperience) Delete(ctx context.Context, workExperienceId int) error {
	return repo.db.TeacherWorkExperience.
		DeleteOneID(workExperienceId).
		Exec(ctx)
}

func (repo *teacherWorkExperience) CreateMany(ctx context.Context, value []*ent.TeacherWorkExperience, teacherId int) error {
	bulk := make([]*ent.TeacherWorkExperienceCreate, len(value))
	for i, v := range value {
		bulk[i] = repo.db.TeacherWorkExperience.Create().
			SetAcademyName(v.AcademyName).
			SetWorkStartAt(v.WorkStartAt).
			SetNillableWorkEndAt(v.WorkEndAt).
			SetNillableDescription(v.Description).
			SetTeacherID(teacherId)
	}

	return repo.db.TeacherWorkExperience.CreateBulk(bulk...).Exec(ctx)
}
