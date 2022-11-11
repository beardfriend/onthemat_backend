package repository

import (
	"context"

	"onthemat/internal/app/common"
	"onthemat/internal/app/utils"
	"onthemat/pkg/ent"
	"onthemat/pkg/ent/acadmey"
	"onthemat/pkg/ent/predicate"
	"onthemat/pkg/ent/user"
	"onthemat/pkg/entx"
)

type AcademyRepository interface {
	Create(ctx context.Context, academy *ent.Acadmey, userId int) error
	Update(ctx context.Context, academy *ent.Acadmey, userId int) error
	Get(ctx context.Context, userId int) (*ent.Acadmey, error)
	List(ctx context.Context, p *common.ListParams) ([]*ent.Acadmey, error)
	Total(ctx context.Context, p *common.TotalParams) (int, error)
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

func (repo *academyRepository) Total(ctx context.Context, p *common.TotalParams) (int, error) {
	clause := repo.db.Debug().Acadmey.Query()
	clause = repo.conditionQuery(ctx, p, clause)
	return clause.Count(ctx)
}

func (repo *academyRepository) List(ctx context.Context, p *common.ListParams) ([]*ent.Acadmey, error) {
	pagination := utils.NewPagination(p.PageNo, p.PageSize)

	clause := repo.db.Debug().Acadmey.
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
			acadmey.FieldCreatedAt,
			acadmey.FieldUpdatedAt,
		).
		Limit(pagination.GetLimit()).
		Offset(pagination.GetOffset())

	useableOrderCol := map[string]string{
		"ID":        acadmey.FieldID,
		"CREATEDAT": acadmey.FieldCreatedAt,
	}

	if p.OrderCol != nil && p.OrderType != nil {
		orderCol := useableOrderCol[*p.OrderCol]

		if *p.OrderType == "ASC" {
			clause.Order(ent.Asc(orderCol))
		} else {
			clause.Order(ent.Desc(orderCol))
		}
	} else {
		clause.Order(ent.Desc(acadmey.FieldID))
	}

	clause = repo.conditionQuery(ctx, &common.TotalParams{
		SearchKey:   p.SearchKey,
		SearchValue: p.SearchValue,
	}, clause)

	return clause.All(ctx)
}

func (repo *academyRepository) conditionQuery(ctx context.Context, p *common.TotalParams, clause *ent.AcadmeyQuery) *ent.AcadmeyQuery {
	useableSearchCol := map[string]func(v string) predicate.Acadmey{
		"NAME": acadmey.NameContains,
		"GU":   acadmey.AddressGuContains,
	}

	if p.SearchKey != nil && p.SearchValue != nil {
		whereFunc := useableSearchCol[*p.SearchKey]
		clause.Where(whereFunc(*p.SearchValue))
	}
	return clause
}
