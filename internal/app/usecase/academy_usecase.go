package usecase

import (
	"context"
	"strings"

	"onthemat/internal/app/common"
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
	Update(ctx context.Context, a *transport.AcademyUpdateRequestBody, userId int) error
	List(ctx context.Context, a *request.AcademyListQueries) ([]*ent.Academy, *utils.PagenationInfo, error)
}

type academyUseCase struct {
	userRepo    repository.UserRepository
	academyRepo repository.AcademyRepository
	areaRepo    repository.AreaRepository
	academySvc  service.AcademyService
}

func NewAcademyUsecase(
	academyRepo repository.AcademyRepository,
	academyService service.AcademyService,
	userRepo repository.UserRepository,
	areaRepo repository.AreaRepository,
) *academyUseCase {
	return &academyUseCase{
		academyRepo: academyRepo,
		academySvc:  academyService,
		userRepo:    userRepo,
		areaRepo:    areaRepo,
	}
}

func (u *academyUseCase) Create(ctx context.Context, academy *transport.AcademyCreateRequestBody, userId int) (err error) {
	if err = u.academySvc.VerifyBusinessMan(academy.BusinessCode); err != nil {
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

	sigungu, err := u.areaRepo.GetSigunGu(ctx, academy.AddressSigungu)
	if err != nil {
		if ent.IsNotFound(err) {
			err = ex.NewNotFoundError(ex.ErrAreaNotFound, nil)
			return
		}
		return
	}

	return u.academyRepo.Create(ctx, &ent.Academy{
		Name:          academy.Name,
		BusinessCode:  academy.BusinessCode,
		CallNumber:    academy.CallNumber,
		AddressRoad:   academy.AddressRoad,
		AddressDetail: academy.AddressDetail,
		Edges: ent.AcademyEdges{
			Sigungu: &ent.AreaSiGungu{ID: sigungu.ID},
		},
	}, userId)
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

func (u *academyUseCase) Update(ctx context.Context, a *transport.AcademyUpdateRequestBody, userId int) (err error) {
	sigungu, err := u.areaRepo.GetSigunGu(ctx, a.AddressSigungu)
	if err != nil {
		if ent.IsNotFound(err) {
			err = ex.NewNotFoundError(ex.ErrAreaNotFound, nil)
			return
		}
		return
	}
	if err = u.academyRepo.Update(ctx, &ent.Academy{
		Name:          a.Name,
		CallNumber:    a.CallNumber,
		AddressRoad:   a.AddressRoad,
		AddressDetail: a.AddressDetail,
		Edges: ent.AcademyEdges{
			Sigungu: &ent.AreaSiGungu{
				ID: sigungu.ID,
			},
		},
	}, userId); err != nil {
		if ent.IsNotFound(err) {
			err = ex.NewNotFoundError(ex.ErrAcademyNotFound, nil)
			return
		}
		return
	}

	return
}

func (u *academyUseCase) List(ctx context.Context, a *request.AcademyListQueries) (result []*ent.Academy, paginationInfo *utils.PagenationInfo, err error) {
	if a.OrderCol != nil {
		*a.OrderCol = strings.ToUpper(*a.OrderCol)
	}

	// if a.SearchKey != nil {
	// 	*a.SearchKey = strings.ToUpper(*a.SearchKey)
	// }

	pginationModule := utils.NewPagination(a.PageNo, a.PageSize)

	total, err := u.academyRepo.Total(ctx, &common.TotalParams{SearchKey: a.SearchKey, SearchValue: a.SearchValue})
	if err != nil {
		if err.Error() == repository.ErrOrderColumnInvalid || err.Error() == repository.ErrSearchColumnInvalid {
			err = ex.NewBadRequestError(ex.ErrColumnInvalid, nil)
			return
		}
		return
	}

	pginationModule.SetTotal(total)
	result, err = u.academyRepo.List(ctx, &common.ListParams{
		PageNo:      a.PageNo,
		PageSize:    a.PageSize,
		OrderCol:    a.OrderCol,
		OrderType:   a.OrderType,
		SearchKey:   a.SearchKey,
		SearchValue: a.SearchValue,
	})
	if err != nil {
		return
	}
	paginationInfo = pginationModule.GetInfo(len(result))
	return
}
