package repository

import (
	"context"

	"onthemat/pkg/ent"
	ri "onthemat/pkg/ent/recruitmentinstead"
)

type recruitmentInsteadRepo struct{}

func (repo *recruitmentInsteadRepo) createMany(ctx context.Context, db *ent.Client, vals []*ent.RecruitmentInstead, recruitmentId int) (err error) {
	bulk := make([]*ent.RecruitmentInsteadCreate, len(vals))
	for i, v := range vals {
		clause := db.RecruitmentInstead.Create().
			SetRecuritmentID(recruitmentId).
			SetMinCareer(v.MinCareer).
			SetPay(v.Pay).
			SetStartDateTime(v.StartDateTime).
			SetEndDateTime(v.EndDateTime)

		if v.ID != 0 {
			clause.SetID(v.ID)
		}

		bulk[i] = clause
	}

	return db.RecruitmentInstead.CreateBulk(bulk...).Exec(ctx)
}

func (repo *recruitmentInsteadRepo) getIdsByRecruitId(ctx context.Context, db *ent.Client, recruitId int) ([]int, error) {
	return db.RecruitmentInstead.Query().
		Where(ri.RecruitmentIDEQ(recruitId)).IDs(ctx)
}

func (repo *recruitmentInsteadRepo) updateMany(ctx context.Context, db *ent.Client, vals []*ent.RecruitmentInstead, recruitId int) (err error) {
	for _, v := range vals {
		err = db.RecruitmentInstead.Update().
			Where(ri.RecruitmentIDEQ(recruitId)).
			SetRecuritmentID(recruitId).
			SetMinCareer(v.MinCareer).
			SetPay(v.Pay).
			SetStartDateTime(v.StartDateTime).
			SetEndDateTime(v.EndDateTime).
			SetNillablePasserID(v.TeacherID).Exec(ctx)
		if err != nil {
			return
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
