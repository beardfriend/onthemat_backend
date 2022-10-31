package repository

import (
	"context"

	"onthemat/pkg/ent"
	"onthemat/pkg/ent/acadmey"
)

type AcademyRepository interface {
	Create(ctx context.Context, academy *ent.Acadmey, userId int) error
	Get(ctx context.Context, userId int) (*ent.Acadmey, error)
}

type academyRepository struct {
	db *ent.Client
}

func NewAcademyRepository(db *ent.Client) AcademyRepository {
	return &academyRepository{
		db: db,
	}
}

func (svc *academyRepository) Create(ctx context.Context, academy *ent.Acadmey, userId int) error {
	return svc.db.Acadmey.Create().
		SetName(academy.Name).
		SetBusinessCode(academy.BusinessCode).
		SetCallNumber(academy.CallNumber).
		SetAddressRoad(academy.AddressRoad).
		SetAddressSigun(academy.AddressSigun).
		SetAddressGu(academy.AddressGu).
		SetAddressDong(academy.AddressDong).
		SetAddressX(academy.AddressX).
		SetAddressY(academy.AddressY).SetUserID(userId).Exec(ctx)
}

func (svc *academyRepository) Get(ctx context.Context, userId int) (*ent.Acadmey, error) {
	return svc.db.Acadmey.
		Query().
		Where(acadmey.ID(userId)).
		WithUser().
		Only(ctx)
}
