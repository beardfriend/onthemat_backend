package repository

import (
	"context"

	"onthemat/pkg/ent"
	tcf "onthemat/pkg/ent/teachercertification"
)

type teacherCertification struct{}

func (repo *teacherCertification) createMany(ctx context.Context, db *ent.Client, value []*ent.TeacherCertification, teacherId int) error {
	bulk := make([]*ent.TeacherCertificationCreate, len(value))
	for i, v := range value {
		clause := db.TeacherCertification.Create().
			SetTeacherID(teacherId).
			SetAgencyName(v.AgencyName).
			SetClassStartAt(v.ClassStartAt).
			SetNillableImageUrl(v.ImageUrl).
			SetNillableClassEndAt(v.ClassEndAt).
			SetNillableDescription(v.Description)

		if v.ID != 0 {
			clause.SetID(v.ID)
		}
		bulk[i] = clause

	}
	return db.TeacherCertification.CreateBulk(bulk...).Exec(ctx)
}

func (repo *teacherCertification) getIdsByTeacherId(ctx context.Context, db *ent.Client, teacherId int) ([]int, error) {
	return db.TeacherCertification.Query().
		Where(tcf.TeacherIDEQ(teacherId)).
		IDs(ctx)
}

func (repo *teacherCertification) updateMany(ctx context.Context, db *ent.Client, value []*ent.TeacherCertification, teacherId int) (err error) {
	for _, v := range value {
		clause := db.TeacherCertification.Update().
			Where(tcf.IDEQ(v.ID), tcf.TeacherIDEQ(teacherId)).
			SetTeacherID(teacherId).
			SetAgencyName(v.AgencyName).
			SetClassStartAt(v.ClassStartAt).
			SetNillableImageUrl(v.ImageUrl).
			SetNillableClassEndAt(v.ClassEndAt).
			SetNillableDescription(v.Description)

		mu := clause.Mutation()
		if v.ImageUrl == nil {
			mu.ClearImageUrl()
		}
		if v.ClassEndAt == nil {
			mu.ClearClassEndAt()
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

func (repo *teacherCertification) deletesByIds(ctx context.Context, db *ent.Client, ids []int, teacherId int) (int, error) {
	return db.TeacherCertification.Delete().Where(tcf.And(
		tcf.TeacherIDEQ(teacherId),
		tcf.IDIn(ids...),
	)).Exec(ctx)
}

func (repo *teacherCertification) deletebyTeacherId(ctx context.Context, db *ent.Client, teacherId int) error {
	_, err := db.TeacherCertification.Delete().Where(tcf.TeacherIDEQ(teacherId)).Exec(ctx)
	return err
}
