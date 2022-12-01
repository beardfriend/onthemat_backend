package usecase

import (
	"context"
	"strings"

	ex "onthemat/internal/app/common"
	"onthemat/internal/app/repository"
	"onthemat/internal/app/service"
	"onthemat/internal/app/transport/request"

	"onthemat/internal/app/utils"
	"onthemat/pkg/ent"
)

type AcademyUsecase interface {
	Create(ctx context.Context, academy *request.AcademyCreateBody, userId int) error
	Put(ctx context.Context, req *request.AcademyUpdateBody, id, userId int) (isUpdated bool, err error)
	Patch(ctx context.Context, req *request.AcademyPatchBody, id, userId int) (err error)
	Get(ctx context.Context, userId int) (*ent.Academy, error)
	List(ctx context.Context, a *request.AcademyListQueries) ([]*ent.Academy, *utils.PagenationInfo, error)
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

func (u *academyUseCase) Create(ctx context.Context, req *request.AcademyCreateBody, userId int) (err error) {
	info := req.Info

	// BusinessCode Check
	if err = u.academySvc.VerifyBusinessMan(info.BusinessCode); err != nil {
		if err.Error() == service.ErrBussinessCodeInvalid {
			err = ex.NewBadRequestError(ex.ErrBusinessCodeInvalid, nil)
			return
		}
		return
	}

	// Already Exisit
	getUser, err := u.userRepo.Get(ctx, userId)
	if err != nil {
		if ent.IsNotFound(err) {
			err = ex.NewNotFoundError(ex.ErrUserNotFound, nil)
			return
		}
		return
	}

	// Check User Type
	if getUser.Type != nil {
		err = ex.NewConflictError(ex.ErrUserTypeAlreadyRegisted, nil)
		return
	}

	// Prepare Data
	yoga := make([]*ent.Yoga, 0)
	if req.YogaIds != nil {
		for _, v := range *req.YogaIds {
			yoga = append(yoga, &ent.Yoga{ID: v})
		}
	}
	data := &ent.Academy{
		UserID:        userId,
		SigunguID:     info.SigunguID,
		Name:          info.Name,
		BusinessCode:  info.BusinessCode,
		CallNumber:    info.CallNumber,
		AddressRoad:   info.AddressRoad,
		AddressDetail: info.AddressDetail,
		Edges: ent.AcademyEdges{
			Yoga: yoga,
		},
	}

	// Do
	err = u.academyRepo.Create(ctx, data)

	if err != nil {
		if ent.IsConstraintError(err) {
			err = foreignKeyConstraintError(err)
			return
		}
		return
	}
	return
}

func foreignKeyConstraintError(err error) error {
	if strings.Contains(err.Error(), "yoga_id") {
		err = ex.NewConflictError(ex.ErrYogaDoseNotExist, nil)
	} else if strings.Contains(err.Error(), "sigungu") {
		err = ex.NewConflictError(ex.ErrSigunguDoseNotExist, nil)
	} else {
		err = ex.NewConflictError(ex.ErrConflict, nil)
	}
	return err
}

func (u *academyUseCase) Put(ctx context.Context, req *request.AcademyUpdateBody, id, userId int) (isUpdated bool, err error) {
	// Prepare Data
	info := req.Info
	yoga := make([]*ent.Yoga, 0)
	if req.YogaIDs != nil {
		for _, v := range *req.YogaIDs {
			yoga = append(yoga, &ent.Yoga{ID: v})
		}
	}
	data := &ent.Academy{
		ID:            id,
		UserID:        userId,
		SigunguID:     info.SigunguID,
		Name:          info.Name,
		CallNumber:    info.CallNumber,
		AddressRoad:   info.AddressRoad,
		AddressDetail: info.AddressDetail,
		Edges: ent.AcademyEdges{
			Yoga: yoga,
		},
	}

	// Already Exisit
	isExist, err := u.academyRepo.Exist(ctx, id)
	isUpdated = !isExist
	if err != nil {
		return
	}

	// Do
	if !isExist {
		err = u.academyRepo.Create(ctx, data)
	} else {
		err = u.academyRepo.Update(ctx, data)
	}

	if err != nil {
		if ent.IsConstraintError(err) {
			err = foreignKeyConstraintError(err)
			return
		}
		return
	}
	return
}

func (u *academyUseCase) Patch(ctx context.Context, req *request.AcademyPatchBody, id, userId int) (err error) {
	err = u.academyRepo.Patch(ctx, req, id, userId)
	if err != nil {
		if ent.IsConstraintError(err) {
			err = foreignKeyConstraintError(err)
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
		return
	}
	orderType := ex.DESC
	if a.OrderType != nil && *a.OrderType == string(ex.ASC) {
		orderType = ex.ASC
	}

	pginationModule.SetTotal(total)
	result, err = u.academyRepo.List(ctx, paginationModule, a.YogaIDs, a.SigunGuID, a.AcademyName, a.OrderCol, orderType)
	if err != nil {
		return
	}

	paginationInfo = pginationModule.GetInfo(len(result))
	return
}
