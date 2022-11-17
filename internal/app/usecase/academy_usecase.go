package usecase

import (
	"context"
	"strings"

	"onthemat/internal/app/common"
	ex "onthemat/internal/app/common"
	"onthemat/internal/app/repository"
	"onthemat/internal/app/service"
	"onthemat/internal/app/transport"
	"onthemat/internal/app/utils"
	"onthemat/pkg/ent"
)

type AcademyUsecase interface {
	Create(ctx context.Context, academy *transport.AcademyCreateRequestBody, userId int) error
	Get(ctx context.Context, userId int) (*ent.Academy, error)
	Update(ctx context.Context, a *transport.AcademyUpdateRequestBody, userId int) error
	List(ctx context.Context, a *transport.AcademyListQueries) ([]*ent.Academy, *utils.PagenationInfo, error)
}

type academyUseCase struct {
	userRepo    repository.UserRepository
	academyRepo repository.AcademyRepository
	academySvc  service.AcademyService
}

func NewAcademyUsecase(
	academyRepo repository.AcademyRepository,
	academyService service.AcademyService,
	userRepo repository.UserRepository,
) *academyUseCase {
	return &academyUseCase{
		academyRepo: academyRepo,
		academySvc:  academyService,
		userRepo:    userRepo,
	}
}

func (u *academyUseCase) Create(ctx context.Context, academy *transport.AcademyCreateRequestBody, userId int) error {
	if err := u.academySvc.VerifyBusinessMan(academy.BusinessCode); err != nil {
		if err.Error() == service.ErrBussinessCodeInvalid {
			return ex.NewBadRequestError(ex.ErrBusinessCodeInvalid, nil)
		}
		return err
	}

	getUser, err := u.userRepo.Get(ctx, userId)
	if err != nil {
		if ent.IsNotFound(err) {
			return ex.NewNotFoundError(ex.ErrUserNotFound, nil)
		}
		return err
	}

	if getUser.Type != nil {
		return ex.NewConflictError(ex.ErrUserTypeAlreadyRegisted, nil)
	}

	return u.academyRepo.Create(ctx, &ent.Academy{
		Name:          academy.Name,
		BusinessCode:  academy.BusinessCode,
		CallNumber:    academy.CallNumber,
		AddressRoad:   academy.AddressRoad,
		AddressSigun:  academy.AddressSigun,
		AddressGu:     academy.AddressGu,
		AddressDong:   academy.AddressDong,
		AddressDetail: academy.AddressDetail,
		AddressX:      academy.AddressX,
		AddressY:      academy.AddressY,
	}, userId)
}

func (u *academyUseCase) Get(ctx context.Context, userId int) (*ent.Academy, error) {
	academy, err := u.academyRepo.Get(ctx, userId)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ex.NewNotFoundError(ex.ErrAcademyNotFound, nil)
		}
	}
	return academy, nil
}

func (u *academyUseCase) Update(ctx context.Context, a *transport.AcademyUpdateRequestBody, userId int) error {
	if err := u.academyRepo.Update(ctx, &ent.Academy{
		Name:          a.Name,
		CallNumber:    a.CallNumber,
		AddressRoad:   a.AddressRoad,
		AddressSigun:  a.AddressSigun,
		AddressGu:     a.AddressGu,
		AddressDong:   a.AddressDong,
		AddressDetail: a.AddressDetail,
		AddressX:      a.AddressX,
		AddressY:      a.AddressY,
	}, userId); err != nil {
		if ent.IsNotFound(err) {
			return ex.NewNotFoundError(ex.ErrAcademyNotFound, nil)
		}
	}

	return nil
}

func (u *academyUseCase) List(ctx context.Context, a *transport.AcademyListQueries) ([]*ent.Academy, *utils.PagenationInfo, error) {
	if a.OrderCol != nil {
		*a.OrderCol = strings.ToUpper(*a.OrderCol)
	}

	pagination := utils.NewPagination(a.PageNo, a.PageSize)

	total, err := u.academyRepo.Total(ctx, &common.TotalParams{SearchKey: a.SearchKey, SearchValue: a.SearchValue})
	if err != nil {
		return nil, nil, err
	}

	pagination.SetTotal(total)
	academis, err := u.academyRepo.List(ctx, &common.ListParams{
		PageNo:      a.PageNo,
		PageSize:    a.PageSize,
		OrderCol:    a.OrderCol,
		OrderType:   a.OrderType,
		SearchKey:   a.SearchKey,
		SearchValue: a.SearchValue,
	})
	if err != nil {
		return nil, nil, err
	}

	return academis, pagination.GetInfo(len(academis)), nil
}
