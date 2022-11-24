package usecase

import (
	"context"

	ex "onthemat/internal/app/common"
	"onthemat/internal/app/repository"
	r "onthemat/internal/app/repository"
	"onthemat/internal/app/transport/request"
	"onthemat/internal/app/utils"
	"onthemat/pkg/ent"
)

type YogaUsecase interface {
	CreateGroup(ctx context.Context, req *request.YogaGroupCreateBody) (err error)
	GroupList(ctx context.Context, req *request.YogaGroupListQueries) (result []*ent.YogaGroup, pagination *utils.PagenationInfo, err error)
	UpdateGroup(ctx context.Context, req *request.YogaGroupUpdateBody, yogaId int) error
	DeleteGroup(ctx context.Context, ids []int) (rowAffected int, err error)

	Create(ctx context.Context, req *request.YogaCreateBody) (err error)
	List(ctx context.Context, groupId int) ([]*ent.Yoga, error)
	Update(ctx context.Context, req *request.YogaUpdateBody, yogaId int) (err error)
	Delete(ctx context.Context, yogaId int) error
}

type yogaUseCase struct {
	yogaRepo r.YogaRepository
}

func NewYogaUsecase(yogaRepo r.YogaRepository) YogaUsecase {
	return &yogaUseCase{
		yogaRepo: yogaRepo,
	}
}

// ------------------- Group -------------------

func (u *yogaUseCase) CreateGroup(ctx context.Context, req *request.YogaGroupCreateBody) (err error) {
	err = u.yogaRepo.CreateGroup(ctx, &ent.YogaGroup{
		Category:    req.Category,
		CategoryEng: req.CategoryEng,
		Description: req.Description,
	})
	if err != nil {
		if ent.IsConstraintError(err) {
			err = ex.NewConflictError(ex.ErrYogaGroupAlreadyExist, nil)
			return
		}
		return
	}
	return
}

func (u *yogaUseCase) DeleteGroup(ctx context.Context, ids []int) (rowAffected int, err error) {
	return u.yogaRepo.DeleteGroups(ctx, ids)
}

func (u *yogaUseCase) UpdateGroup(ctx context.Context, req *request.YogaGroupUpdateBody, yogaId int) error {
	return u.yogaRepo.UpdateGroup(ctx, &ent.YogaGroup{
		ID:          yogaId,
		Category:    req.Category,
		CategoryEng: req.CategoryEng,
		Description: req.Description,
	})
}

func (u *yogaUseCase) GroupList(ctx context.Context, req *request.YogaGroupListQueries) (result []*ent.YogaGroup, pagination *utils.PagenationInfo, err error) {
	pgModule := utils.NewPagination(req.PageNo, req.PageSize)

	total, err := u.yogaRepo.GroupTotal(ctx, &ex.TotalParams{SearchKey: req.SearchKey, SearchValue: req.SearchValue})
	if err != nil {
		if err.Error() == repository.ErrOrderColumnInvalid || err.Error() == repository.ErrSearchColumnInvalid {
			err = ex.NewBadRequestError(ex.ErrColumnInvalid, nil)
			return
		}
		return
	}
	pgModule.SetTotal(total)
	result, err = u.yogaRepo.GroupList(ctx, pgModule, &ex.ListParams{
		SearchKey:   req.SearchKey,
		SearchValue: req.SearchValue,
	})
	if err != nil {
		return
	}

	pagination = pgModule.GetInfo(len(result))
	return
}

// ------------------- Yoga -------------------

func (u *yogaUseCase) Create(ctx context.Context, req *request.YogaCreateBody) (err error) {
	return u.yogaRepo.Create(ctx, &ent.Yoga{
		NameKor:     req.NameKor,
		NameEng:     req.NameEng,
		Description: req.Description,
		Level:       req.Level,
		Edges: ent.YogaEdges{
			YogaGroup: &ent.YogaGroup{
				ID: req.YogaGroupId,
			},
		},
	})
}

func (u *yogaUseCase) Update(ctx context.Context, req *request.YogaUpdateBody, yogaId int) (err error) {
	return u.yogaRepo.Update(ctx, &ent.Yoga{
		ID:          yogaId,
		NameKor:     req.NameKor,
		NameEng:     req.NameEng,
		Description: req.Description,
		Level:       req.Level,
		Edges: ent.YogaEdges{
			YogaGroup: &ent.YogaGroup{
				ID: req.YogaGroupId,
			},
		},
	})
}

func (u *yogaUseCase) Delete(ctx context.Context, yogaId int) error {
	return u.yogaRepo.Delete(ctx, yogaId)
}

func (u *yogaUseCase) List(ctx context.Context, groupId int) ([]*ent.Yoga, error) {
	return u.yogaRepo.List(ctx, groupId)
}
