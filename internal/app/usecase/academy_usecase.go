package usecase

import (
	"context"

	"onthemat/internal/app/repository"
	"onthemat/internal/app/transport"
	"onthemat/pkg/ent"
)

type AcademyUsecase interface{}

type academyUseCase struct {
	academyRepo repository.AcademyRepository
}

func NewAcademyUsecase(academyRepo repository.AcademyRepository) *academyUseCase {
	return &academyUseCase{
		academyRepo: academyRepo,
	}
}

func (u *academyUseCase) Create(ctx context.Context, academy *transport.AcademyCreateRequestBody, userId int) error {
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

func (u *academyUseCase) Get(ctx context.Context, userId int) {
	u.academyRepo.Get(ctx, userId)
}
