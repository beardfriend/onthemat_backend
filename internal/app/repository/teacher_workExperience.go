package repository

import (
	"context"

	"onthemat/pkg/ent"
	twe "onthemat/pkg/ent/teacherworkexperience"
)

type teacherWorkExperience struct{}

func (repo *teacherWorkExperience) createMany(ctx context.Context, db *ent.Client, value []*ent.TeacherWorkExperience, teacherId int) error {
	bulk := make([]*ent.TeacherWorkExperienceCreate, len(value))
	for i, v := range value {
		clause := db.TeacherWorkExperience.Create().
			SetAcademyName(v.AcademyName).
			SetWorkStartAt(v.WorkStartAt).
			SetNillableWorkEndAt(v.WorkEndAt).
			SetNillableDescription(v.Description).
			SetTeacherID(teacherId)

		if v.ID != 0 {
			clause.SetID(v.ID)
		}
		bulk[i] = clause
	}

	return db.TeacherWorkExperience.CreateBulk(bulk...).Exec(ctx)
}

func (repo *teacherWorkExperience) getIdsByTeacherId(ctx context.Context, db *ent.Client, teacherId int) ([]int, error) {
	return db.TeacherWorkExperience.Query().
		Where(twe.TeacherIDEQ(teacherId)).
		IDs(ctx)
}

func (repo *teacherWorkExperience) updateMany(ctx context.Context, db *ent.Client, value []*ent.TeacherWorkExperience, teacherId int) (err error) {
	for _, v := range value {
		clause := db.TeacherWorkExperience.Update().
			Where(twe.IDEQ(v.ID), twe.TeacherIDEQ(teacherId)).
			SetTeacherID(teacherId).
			SetAcademyName(v.AcademyName).
			SetWorkStartAt(v.WorkStartAt).
			SetNillableWorkEndAt(v.WorkEndAt).
			SetNillableDescription(v.Description)

		mu := clause.Mutation()
		if v.WorkEndAt == nil {
			mu.ClearWorkEndAt()
		}
		if v.Description == nil {
			mu.ClearDescription()
		}

		err = clause.Exec(ctx)
		if err != nil {
			return err
		}
	}
	return
}

func (repo *teacherWorkExperience) deletesByIds(ctx context.Context, db *ent.Client, ids []int, teacherId int) (int, error) {
	return db.TeacherWorkExperience.Delete().Where(twe.And(
		twe.TeacherIDEQ(teacherId),
		twe.IDIn(ids...),
	)).Exec(ctx)
}

func (repo *teacherWorkExperience) deletesByTecaherId(ctx context.Context, db *ent.Client, teacherId int) error {
	_, err := db.TeacherWorkExperience.Delete().Where(twe.TeacherIDEQ(teacherId)).Exec(ctx)
	return err
}
