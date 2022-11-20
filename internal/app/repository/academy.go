package repository

import (
	"context"
	"errors"

	"onthemat/internal/app/common"
	"onthemat/internal/app/model"
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

const (
	ErrSearchColumnInvalid = "유효하지 않은 Column입니다"
	ErrOrderColumnInvalid  = "유효하지 않은 Column입니다"
)

func NewAcademyRepository(db *ent.Client) AcademyRepository {
	return &academyRepository{
		db: db,
	}
}

func (repo *academyRepository) Create(ctx context.Context, academy *ent.Academy, userId int) error {
	return entx.WithTx(ctx, repo.db, func(tx *ent.Tx) (err error) {
		if err = repo.db.Academy.Create().
			SetName(academy.Name).
			SetBusinessCode(academy.BusinessCode).
			SetCallNumber(academy.CallNumber).
			SetAddressRoad(academy.AddressRoad).
			SetAddressSigun(academy.AddressSigun).
			SetAddressGu(academy.AddressGu).
			SetAddressDong(academy.AddressDong).
			SetAddressX(academy.AddressX).
			SetAddressY(academy.AddressY).SetUserID(userId).Exec(ctx); err != nil {
			return
		}

		if err = repo.db.User.UpdateOneID(userId).SetType(model.AcademyType).Exec(ctx); err != nil {
			return
		}

		return
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
			academy.FieldCreatedAt,
			academy.FieldUpdatedAt,
		).
		Where(academy.ID(userId)).
		Only(ctx)
}

func (repo *academyRepository) Total(ctx context.Context, p *common.TotalParams) (total int, err error) {
	clause := repo.db.Academy.Query()
	clause, err = repo.conditionQuery(ctx, p, clause)
	if err != nil {
		return
	}
	total, err = clause.Count(ctx)
	return
}

func (repo *academyRepository) List(ctx context.Context, p *common.ListParams) (result []*ent.Academy, err error) {
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

	if p.OrderCol == nil || p.OrderType == nil {
		clause.Order(ent.Desc(academy.FieldID))
	}

	if p.OrderCol != nil && p.OrderType != nil {
		orderCol := useableOrderCol[*p.OrderCol]

		if orderCol == "" {
			err = errors.New(ErrOrderColumnInvalid)
			return
		}

		if *p.OrderType == "ASC" {
			clause.Order(ent.Asc(orderCol))
		} else {
			clause.Order(ent.Desc(orderCol))
		}

	}

	clause, err = repo.conditionQuery(ctx, &common.TotalParams{
		SearchKey:   p.SearchKey,
		SearchValue: p.SearchValue,
	}, clause)

	if err != nil {
		return
	}
	result, err = clause.All(ctx)
	return
}

func (repo *academyRepository) conditionQuery(
	ctx context.Context,
	p *common.TotalParams,
	clause *ent.AcademyQuery,
) (*ent.AcademyQuery, error) {
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
