package usecase

import (
	"context"
	"fmt"
	"strings"

	ex "onthemat/internal/app/common"
	"onthemat/internal/app/repository"
	"onthemat/internal/app/service"
	"onthemat/internal/app/transport"
	"onthemat/internal/app/transport/request"

	"onthemat/internal/app/utils"
	"onthemat/pkg/ent"
)

type AcademyUsecase interface {
	Create(ctx context.Context, academy *transport.AcademyCreateRequestBody, userId int) error
	Get(ctx context.Context, userId int) (*ent.Academy, error)
	List(ctx context.Context, a *request.AcademyListQueries) ([]*ent.Academy, *utils.PagenationInfo, error)
	Patch(ctx context.Context, req *transport.AcademyPatchRequestBody, id int) error
}

type academyUseCase struct {
	userRepo    repository.UserRepository
	yogaRepo    repository.YogaRepository
	academyRepo repository.AcademyRepository
	areaRepo    repository.AreaRepository
	academySvc  service.AcademyService
}

func NewAcademyUsecase(
	academyRepo repository.AcademyRepository,
	academyService service.AcademyService,
	userRepo repository.UserRepository,
	yogaRepo repository.YogaRepository,
	areaRepo repository.AreaRepository,
) AcademyUsecase {
	return &academyUseCase{
		academyRepo: academyRepo,
		academySvc:  academyService,
		userRepo:    userRepo,
		areaRepo:    areaRepo,
		yogaRepo:    yogaRepo,
	}
}

func (u *academyUseCase) Create(ctx context.Context, req *transport.AcademyCreateRequestBody, userId int) (err error) {
	info := req.Info
	if err = u.academySvc.VerifyBusinessMan(info.BusinessCode); err != nil {
		if err.Error() == service.ErrBussinessCodeInvalid {
			err = ex.NewBadRequestError(ex.ErrBusinessCodeInvalid, nil)
			return
		}
		return
	}

	getUser, err := u.userRepo.Get(ctx, userId)
	if err != nil {
		if ent.IsNotFound(err) {
			err = ex.NewNotFoundError(ex.ErrUserNotFound, nil)
			return
		}
		return
	}

	if getUser.Type != nil {
		err = ex.NewConflictError(ex.ErrUserTypeAlreadyRegisted, nil)
		return
	}

	_, err = u.areaRepo.GetSigunGu(ctx, info.AddressSigungu)
	if err != nil {
		if ent.IsNotFound(err) {
			err = ex.NewNotFoundError(ex.ErrAreaNotFound, nil)
			return
		}
		return
	}

	return
}

func (u *academyUseCase) Get(ctx context.Context, userId int) (result *ent.Academy, err error) {
	result, err = u.academyRepo.Get(ctx, userId)
	if err != nil {
		if ent.IsNotFound(err) {
			err = ex.NewNotFoundError(ex.ErrAcademyNotFound, nil)
			return
		}
		return
	}

	return
}

func (u *academyUseCase) List(ctx context.Context, a *request.AcademyListQueries) (result []*ent.Academy, paginationInfo *utils.PagenationInfo, err error) {
	paginationModule := utils.NewPagination(a.PageNo, a.PageSize)

	if a.OrderCol != nil {
		*a.OrderCol = strings.ToUpper(*a.OrderCol)
	}

	pginationModule := utils.NewPagination(a.PageNo, a.PageSize)

	total, err := u.academyRepo.Total(ctx, a.YogaIDs, a.SigunGuID, a.AcademyName)
	if err != nil {
		if err.Error() == repository.ErrOrderColumnInvalid || err.Error() == repository.ErrSearchColumnInvalid {
			err = ex.NewBadRequestError(ex.ErrColumnInvalid, nil)
			return
		}
		return
	}

	pginationModule.SetTotal(total)
	result, err = u.academyRepo.List(ctx, paginationModule, a.YogaIDs, a.SigunGuID, a.AcademyName, a.OrderCol, a.OrderType)
	if err != nil {
		return
	}

	paginationInfo = pginationModule.GetInfo(len(result))
	return
}

func (u *academyUseCase) Patch(ctx context.Context, req *transport.AcademyPatchRequestBody, id int) error {
	aInfo := req.Info

	d := &ent.Academy{
		ID:            id,
		Name:          aInfo.Name,
		CallNumber:    aInfo.CallNumber,
		AddressRoad:   aInfo.AddressRoad,
		AddressDetail: aInfo.AddressDetail,
		SigunguID:     aInfo.SigunguId,
	}

	for _, v := range req.YogaIDs {
		d.Edges.Yoga = append(d.Edges.Yoga, &ent.Yoga{
			ID: v,
		})
	}
	err := u.academyRepo.Patch(ctx, d)
	fmt.Println(err)
	return err
}
