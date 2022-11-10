package repository

import (
	"context"

	"onthemat/internal/app/utils"
	"onthemat/pkg/ent"
	"onthemat/pkg/ent/acadmey"
	"onthemat/pkg/ent/user"
	"onthemat/pkg/entx"
)

type AcademyRepository interface {
	Create(ctx context.Context, academy *ent.Acadmey, userId int) error
	Update(ctx context.Context, academy *ent.Acadmey, userId int) error
	Get(ctx context.Context, userId int) (*ent.Acadmey, error)
	List(ctx context.Context, pageNo, pageSize int) ([]*ent.Acadmey, error)
}

type academyRepository struct {
	db *ent.Client
}

func NewAcademyRepository(db *ent.Client) AcademyRepository {
	return &academyRepository{
		db: db,
	}
}

func (repo *academyRepository) Create(ctx context.Context, academy *ent.Acadmey, userId int) error {
	return entx.WithTx(ctx, repo.db, func(tx *ent.Tx) error {
		if err := repo.db.Acadmey.Create().
			SetName(academy.Name).
			SetBusinessCode(academy.BusinessCode).
			SetCallNumber(academy.CallNumber).
			SetAddressRoad(academy.AddressRoad).
			SetAddressSigun(academy.AddressSigun).
			SetAddressGu(academy.AddressGu).
			SetAddressDong(academy.AddressDong).
			SetAddressX(academy.AddressX).
			SetAddressY(academy.AddressY).SetUserID(userId).Exec(ctx); err != nil {
			return err
		}

		if err := repo.db.User.UpdateOneID(userId).SetType(user.TypeAcademy).Exec(ctx); err != nil {
			return err
		}

		return nil
	})
}

func (repo *academyRepository) Update(ctx context.Context, academy *ent.Acadmey, userId int) error {
	return repo.db.Debug().Acadmey.Update().
		SetName(academy.Name).
		SetCallNumber(academy.CallNumber).
		SetAddressRoad(academy.AddressRoad).
		SetAddressSigun(academy.AddressSigun).
		SetAddressGu(academy.AddressGu).
		SetAddressDong(academy.AddressDong).
		SetAddressDong(academy.AddressDetail).
		SetAddressX(academy.AddressX).
		SetAddressY(academy.AddressY).
		Where(acadmey.IDEQ(userId)).
		Exec(ctx)
}

func (repo *academyRepository) Get(ctx context.Context, userId int) (*ent.Acadmey, error) {
	return repo.db.Acadmey.
		Query().
		Select(
			acadmey.FieldID,
			acadmey.FieldCallNumber,
			acadmey.FieldName,
			acadmey.FieldAddressSigun,
			acadmey.FieldAddressGu,
			acadmey.FieldAddressDong,
			acadmey.FieldAddressRoad,
			acadmey.FieldAddressDetail,
			acadmey.FieldAddressX,
			acadmey.FieldAddressY,
		).
		Where(acadmey.ID(userId)).
		Only(ctx)
}

func (repo *academyRepository) List(ctx context.Context, pageNo, pageSize int) ([]*ent.Acadmey, error) {
	pagination := utils.NewPagination(pageNo, pageSize)

	return repo.db.Acadmey.
		Query().
		Select(
			acadmey.FieldID,
			acadmey.FieldCallNumber,
			acadmey.FieldName,
			acadmey.FieldAddressSigun,
			acadmey.FieldAddressGu,
			acadmey.FieldAddressDong,
			acadmey.FieldAddressRoad,
			acadmey.FieldAddressDetail,
		).
		Limit(pagination.GetLimit()).
		Offset(pagination.GetOffset()).
		All(ctx)
}
