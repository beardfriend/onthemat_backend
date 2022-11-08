package usecase

import (
	"context"

	ex "onthemat/internal/app/common"
	"onthemat/internal/app/repository"
	"onthemat/internal/app/service"
	"onthemat/internal/app/transport"
	"onthemat/pkg/ent"
	"onthemat/pkg/ent/user"
)

type AcademyUsecase interface {
	Create(ctx context.Context, academy *transport.AcademyCreateRequestBody, userId int) error
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
		return ex.NewBadRequestError(err.Error())
	}

	getUser, err := u.userRepo.Get(ctx, userId)
	if err != nil {
		return err
	}

	if getUser.Type != nil {
		return ex.NewConflictError("이미 존재하는 유저입니다.")
	}

	userType := user.TypeAcademy
	_, err = u.userRepo.Update(ctx, &ent.User{
		ID:   userId,
		Type: &userType,
	})
	if err != nil {
		return err
	}

	return u.academyRepo.Create(ctx, &ent.Acadmey{
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

func (u *academyUseCase) Get(ctx context.Context, userId int) (*ent.Acadmey, error) {
	return u.academyRepo.Get(ctx, userId)
}
