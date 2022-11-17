package repository

import (
	"context"
	"errors"

	"onthemat/internal/app/common"
	"onthemat/internal/app/utils"
	"onthemat/pkg/ent"
	"onthemat/pkg/ent/academy"
	"onthemat/pkg/ent/predicate"
	"onthemat/pkg/entx"
)

type AcademyRepository interface {
	Create(ctx context.Context, academy *ent.Academy, userId int) error
	Update(ctx context.Context, academy *ent.Academy, userId int) error
	Get(ctx context.Context, userId int) (*ent.Academy, error)
	List(ctx context.Context, p *common.ListParams) ([]*ent.Academy, error)
	Total(ctx context.Context, p *common.TotalParams) (int, error)
}

type academyRepository struct {
	db *ent.Client
}

const ErrSearchColumnInvalid = "유효하지 않은 Column입니다"

func NewAcademyRepository(db *ent.Client) AcademyRepository {
	return &academyRepository{
		db: db,
	}
}

func (repo *academyRepository) Create(ctx context.Context, academy *ent.Academy, userId int) error {
	return entx.WithTx(ctx, repo.db, func(tx *ent.Tx) error {
		if err := repo.db.Academy.Create().
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

		if err := repo.db.User.UpdateOneID(userId).SetType(1).Exec(ctx); err != nil {
			return err
		}

		return nil
	})
}

func (repo *academyRepository) Update(ctx context.Context, aca *ent.Academy, userId int) error {
	return repo.db.Academy.UpdateOneID(userId).
		SetName(aca.Name).
		SetCallNumber(aca.CallNumber).
		SetAddressRoad(aca.AddressRoad).
		SetAddressSigun(aca.AddressSigun).
		SetAddressGu(aca.AddressGu).
		SetAddressDong(aca.AddressDong).
		SetAddressDong(aca.AddressDetail).
		SetAddressX(aca.AddressX).
		SetAddressY(aca.AddressY).
		Exec(ctx)
}

func (repo *academyRepository) Get(ctx context.Context, userId int) (*ent.Academy, error) {
	return repo.db.Academy.
		Query().
		Select(
			academy.FieldID,
			academy.FieldCallNumber,
			academy.FieldName,
			academy.FieldAddressSigun,
			academy.FieldAddressGu,
			academy.FieldAddressDong,
			academy.FieldAddressRoad,
			academy.FieldAddressDetail,
			academy.FieldAddressX,
			academy.FieldAddressY,
		).
		Where(academy.ID(userId)).
		Only(ctx)
}

func (repo *academyRepository) Total(ctx context.Context, p *common.TotalParams) (int, error) {
	clause := repo.db.Academy.Query()
	clause, err := repo.conditionQuery(ctx, p, clause)
	if err != nil {
		return 0, err
	}
	return clause.Count(ctx)
}

func (repo *academyRepository) List(ctx context.Context, p *common.ListParams) ([]*ent.Academy, error) {
	pagination := utils.NewPagination(p.PageNo, p.PageSize)

	clause := repo.db.Academy.
		Query().
		Select(
			academy.FieldID,
			academy.FieldCallNumber,
			academy.FieldName,
			academy.FieldAddressSigun,
			academy.FieldAddressGu,
			academy.FieldAddressDong,
			academy.FieldAddressRoad,
			academy.FieldAddressDetail,
			academy.FieldCreatedAt,
			academy.FieldUpdatedAt,
		).
		Limit(pagination.GetLimit()).
		Offset(pagination.GetOffset())

	useableOrderCol := map[string]string{
		"ID":        academy.FieldID,
		"CREATEDAT": academy.FieldCreatedAt,
	}

	if p.OrderCol != nil && p.OrderType != nil {
		orderCol := useableOrderCol[*p.OrderCol]

		if orderCol == "" {
			// 에러를 던질 것인지.
			clause.Order(ent.Desc(academy.FieldID))
		}

		if *p.OrderType == "ASC" {
			clause.Order(ent.Asc(orderCol))
		} else {
			clause.Order(ent.Desc(orderCol))
		}
	}

	clause, err := repo.conditionQuery(ctx, &common.TotalParams{
		SearchKey:   p.SearchKey,
		SearchValue: p.SearchValue,
	}, clause)
	if err != nil {
		return nil, err
	}

	return clause.All(ctx)
}

func (repo *academyRepository) conditionQuery(ctx context.Context, p *common.TotalParams, clause *ent.AcademyQuery) (*ent.AcademyQuery, error) {
	useableSearchCol := map[string]func(v string) predicate.Academy{
		"NAME": academy.NameContains,
		"GU":   academy.AddressGuContains,
	}

	if p.SearchKey != nil && p.SearchValue != nil {

		whereFunc := useableSearchCol[*p.SearchKey]

		if whereFunc == nil {
			return nil, errors.New(ErrSearchColumnInvalid)
		}

		clause.Where(whereFunc(*p.SearchValue))
	}
	return clause, nil
}
