package usecase

import (
	"context"

	ex "onthemat/internal/app/common"
	"onthemat/internal/app/model"
	"onthemat/internal/app/repository"
	"onthemat/internal/app/transport/request"
	"onthemat/pkg/ent"
)

type RecruitmentUsecase interface {
	Create(ctx context.Context, d *request.RecruitmentCreateBody, academyId int) (err error)
}

type recruitmentUsecase struct {
	recruitRepo repository.RecruitmentRepository
}

func NewRecruitmentUsecase(recruitRepo repository.RecruitmentRepository) RecruitmentUsecase {
	return &recruitmentUsecase{
		recruitRepo: recruitRepo,
	}
}

func (u *recruitmentUsecase) Create(ctx context.Context, d *request.RecruitmentCreateBody, academyId int) (err error) {
	info := d.Info

	// Prepare Data
	insteadInfo := make([]*ent.RecruitmentInstead, 0)
	for _, v := range d.InsteadInfo {
		schedules := make([]*model.Schedule, 0)

		for _, s := range v.Schedules {
			schedules = append(schedules, &model.Schedule{
				StartDateTime: s.StartDateTime,
				EndDateTime:   s.EndDateTime,
			})
		}
		insteadInfo = append(insteadInfo, &ent.RecruitmentInstead{
			MinCareer: v.MinCareer,
			Pay:       v.Pay,
			Schedule:  schedules,
		})
	}

	data := &ent.Recruitment{
		IsOpen:    info.IsOpen,
		AcademyID: academyId,
		Edges: ent.RecruitmentEdges{
			RecruitmentInstead: insteadInfo,
		},
	}

	err = u.recruitRepo.Create(ctx, data)
	if err != nil {
		if ent.IsConstraintError(err) {
			err = ex.NewConflictError(ex.ErrConflict, nil)
		}
	}
	return
}
