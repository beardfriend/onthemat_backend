package usecase

import (
	"context"

	ex "onthemat/internal/app/common"
	r "onthemat/internal/app/repository"
	"onthemat/internal/app/transport/request"
	"onthemat/pkg/ent"
)

type YogaUsecase interface{}

type yogaUseCase struct {
	yogaRepo r.YogaRepository
}

func NewYogaUsecase(yogaRepo r.YogaRepository) YogaUsecase {
	return &yogaUseCase{
		yogaRepo: yogaRepo,
	}
}

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

func (u *yogaUseCase) CreateYoga(ctx context.Context, req *request.YogaCreateBody) (err error) {
	err = u.yogaRepo.Create(ctx, &ent.Yoga{
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
	if err != nil {
		return
	}
	return
}
