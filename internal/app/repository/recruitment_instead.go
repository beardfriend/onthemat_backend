package repository

import (
	"context"
	"errors"

	"onthemat/internal/app/model"
	"onthemat/pkg/ent"
	ri "onthemat/pkg/ent/recruitmentinstead"
)

type recruitmentInsteadRepo struct{}

func (repo *recruitmentInsteadRepo) createMany(ctx context.Context, db *ent.Client, vals []*ent.RecruitmentInstead, recruitmentId int) (err error) {
	bulk := make([]*ent.RecruitmentInsteadCreate, len(vals))
	for i, v := range vals {

		schedules := make([]*model.Schedule, 0)
		for _, s := range v.Schedule {
			schedules = append(schedules, &model.Schedule{
				StartDateTime: s.StartDateTime,
				EndDateTime:   s.EndDateTime,
			})
		}

		clause := db.RecruitmentInstead.Create().
			SetRecuritmentID(recruitmentId).
			SetMinCareer(v.MinCareer).
			SetPay(v.Pay).
			SetSchedule(schedules)

		if v.ID != 0 {
			clause.SetID(v.ID)
		}

		bulk[i] = clause
	}

	err = db.RecruitmentInstead.CreateBulk(bulk...).Exec(ctx)
	if err != nil {
		if err.Error() == "incosistent id values for batch insert" {
			err = errors.New(ErrOnlyOwnUser)
		}
		return
	}
	return
}

func (repo *recruitmentInsteadRepo) getIdsByRecruitId(ctx context.Context, db *ent.Client, recruitId int) ([]int, error) {
	return db.RecruitmentInstead.Query().
		Where(ri.RecruitmentIDEQ(recruitId)).IDs(ctx)
}

func (repo *recruitmentInsteadRepo) updateMany(ctx context.Context, db *ent.Client, vals []*ent.RecruitmentInstead, recruitId int) (err error) {
	for _, v := range vals {
		schedules := make([]*model.Schedule, 0)
		for _, s := range v.Schedule {
			schedules = append(schedules, &model.Schedule{
				StartDateTime: s.StartDateTime,
				EndDateTime:   s.EndDateTime,
			})
		}

		rowAffcted, err := db.RecruitmentInstead.Update().
			Where(
				ri.IDEQ(v.ID),
				ri.RecruitmentIDEQ(recruitId),
			).
			SetRecuritmentID(recruitId).
			SetMinCareer(v.MinCareer).
			SetPay(v.Pay).
			SetSchedule(schedules).
			SetNillablePasserID(v.TeacherID).Save(ctx)

		if rowAffcted < 1 {
			err = errors.New(ErrOnlyOwnUser)
		}

		if err != nil {
			return err
		}
	}
	return
}

func (repo *recruitmentInsteadRepo) deleteByIds(ctx context.Context, db *ent.Client, ids []int) error {
	_, err := db.RecruitmentInstead.Delete().Where(ri.IDIn(ids...)).Exec(ctx)
	return err
}

func (rep *recruitmentInsteadRepo) deleteByRecruitId(ctx context.Context, db *ent.Client, recruitId int) error {
	_, err := db.RecruitmentInstead.Delete().Where(ri.RecruitmentIDEQ(recruitId)).Exec(ctx)
	return err
}
