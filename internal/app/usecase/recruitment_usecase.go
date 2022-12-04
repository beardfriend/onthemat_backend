package usecase

import (
	"context"

	ex "onthemat/internal/app/common"
	"onthemat/internal/app/model"
	"onthemat/internal/app/repository"
	"onthemat/internal/app/transport/request"
	"onthemat/internal/app/utils"
	"onthemat/pkg/ent"
)

type RecruitmentUsecase interface {
	Create(ctx context.Context, d *request.RecruitmentCreateBody, academyId int) (err error)
	Update(ctx context.Context, d *request.RecruitmentUpdateBody, id, academyId int) (isUpdated bool, err error)
	Patch(ctx context.Context, d *request.RecruitmentPatchBody, id, academyId int) (isUpdated bool, err error)
	SoftDelete(ctx context.Context, id, academyId int) (err error)
	List(ctx context.Context, a *request.RecruitmentListQueries) (result []*ent.Recruitment, paginationInfo *utils.PagenationInfo, err error)
	Get(ctx context.Context, id int) (result *ent.Recruitment, err error)
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

			if err.Error() == repository.ErrOnlyOwnUser {
				err = ex.NewConflictError(ex.ErrResourceUnOwned, nil)
			}
		}
		return
	}
	return
}

func (u *recruitmentUsecase) Update(ctx context.Context, d *request.RecruitmentUpdateBody, id, academyId int) (isUpdated bool, err error) {
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
			ID:        v.Id,
			MinCareer: v.MinCareer,
			Pay:       v.Pay,
			Schedule:  schedules,
		})
	}

	data := &ent.Recruitment{
		ID:        id,
		IsFinish:  info.IsFinish,
		IsOpen:    info.IsOpen,
		AcademyID: academyId,
		Edges: ent.RecruitmentEdges{
			RecruitmentInstead: insteadInfo,
		},
	}

	isExist, err := u.recruitRepo.Exist(ctx, id)
	isUpdated = isExist
	if err != nil {
		return
	}

	if !isExist {
		err = u.recruitRepo.Create(ctx, data)
	} else {
		err = u.recruitRepo.Update(ctx, data)
	}

	if err != nil {
		if err.Error() == repository.ErrOnlyOwnUser {
			err = ex.NewConflictError(ex.ErrResourceUnOwned, nil)
		}

		if ent.IsConstraintError(err) {
			err = ex.NewConflictError(ex.ErrConflict, nil)
		}
		return
	}
	return
}

func (u *recruitmentUsecase) Patch(ctx context.Context, d *request.RecruitmentPatchBody, id, academyId int) (isUpdated bool, err error) {
	isCreated, err := u.recruitRepo.Patch(ctx, d, id, academyId)
	if err != nil {
		if ent.IsConstraintError(err) {
			err = ex.NewConflictError(ex.ErrConflict, nil)
			return
		}
		return
	}
	isUpdated = !isCreated
	return
}

func (u *recruitmentUsecase) List(ctx context.Context, a *request.RecruitmentListQueries) (result []*ent.Recruitment, paginationInfo *utils.PagenationInfo, err error) {
	paginationModule := utils.NewPagination(a.PageNo, a.PageSize)

	pginationModule := utils.NewPagination(a.PageNo, a.PageSize)
	total, err := u.recruitRepo.Total(ctx, a.StartDateTime, a.EndDateTime, a.YogaIDs, a.SigunguIds)
	if err != nil {
		return
	}

	pginationModule.SetTotal(total)
	result, err = u.recruitRepo.List(ctx, paginationModule, a.StartDateTime, a.EndDateTime, a.YogaIDs, a.SigunguIds)
	if err != nil {
		return
	}

	paginationInfo = pginationModule.GetInfo(len(result))
	return
}

func (u *recruitmentUsecase) Get(ctx context.Context, id int) (result *ent.Recruitment, err error) {
	result, err = u.recruitRepo.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			err = ex.NewNotFoundError(ex.ErrRecruitmentNotFound, nil)
			return
		}
		return
	}
	return
}

func (u *recruitmentUsecase) SoftDelete(ctx context.Context, id, academyId int) (err error) {
	err = u.recruitRepo.PatchDeletedAt(ctx, id, academyId)

	if err != nil {
		if err.Error() == repository.ErrOnlyOwnUser {
			err = ex.NewConflictError(ex.ErrResourceUnOwned, nil)
		}

		if ent.IsConstraintError(err) {
			err = ex.NewConflictError(ex.ErrConflict, nil)
		}
		return
	}

	return
}
