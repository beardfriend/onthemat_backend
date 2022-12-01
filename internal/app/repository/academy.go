package repository

import (
	"context"

	"onthemat/internal/app/common"
	"onthemat/internal/app/model"
	"onthemat/internal/app/transport/request"
	"onthemat/internal/app/utils"

	"onthemat/pkg/ent"
	"onthemat/pkg/ent/academy"
	"onthemat/pkg/ent/areasigungu"
	"onthemat/pkg/ent/user"
	"onthemat/pkg/ent/yoga"
	"onthemat/pkg/ent/yogaraw"
	"onthemat/pkg/entx"

	"github.com/fatih/structs"
)

type AcademyRepository interface {
	Create(ctx context.Context, d *ent.Academy) error
	Update(ctx context.Context, d *ent.Academy) error
	Patch(ctx context.Context, d *request.AcademyPatchBody, id, userId int) error
	Exist(ctx context.Context, id int) (bool, error)
	Get(ctx context.Context, id int) (*ent.Academy, error)
	List(ctx context.Context,
		pgModule *utils.Pagination,
		yogaIDs *[]int, sigunguID *int, academyName *string,
		orderCol *string, orderType common.Sorts) (result []*ent.Academy, err error)

	Total(ctx context.Context, yogaIDs *[]int, sigunguID *int, academyName *string) (result int, err error)

	//
	GetOnlyIdByUserId(ctx context.Context, userId int) (id int, err error)
}

type academyRepository struct {
	db *ent.Client
}

func NewAcademyRepository(db *ent.Client) AcademyRepository {
	return &academyRepository{
		db: db,
	}
}

func (repo *academyRepository) Create(ctx context.Context, d *ent.Academy) (err error) {
	return entx.WithTx(ctx, repo.db, func(tx *ent.Tx) (err error) {
		clause := tx.Academy.Create().
			SetUserID(d.UserID).
			SetSigunguID(d.SigunguID).
			SetName(d.Name).
			SetBusinessCode(d.BusinessCode).
			SetCallNumber(d.CallNumber).
			SetAddressRoad(d.AddressRoad).
			SetNillableAddressDetail(d.AddressDetail)

		if len(d.Edges.Yoga) > 0 {
			clause.AddYoga(d.Edges.Yoga...)
		}

		if len(d.Edges.YogaRaw) > 0 {
			clause.AddYogaRaw(d.Edges.YogaRaw...)
		}

		if d.ID != 0 {
			clause.SetID(d.ID)
		}

		if err = clause.Exec(ctx); err != nil {
			return
		}

		err = tx.User.Update().
			SetType(model.AcademyType).
			Where(user.IDEQ(d.UserID)).
			Exec(ctx)
		if err != nil {
			return
		}

		return
	})
}

func (repo *academyRepository) Exist(ctx context.Context, id int) (bool, error) {
	return repo.db.Academy.Query().Where(academy.IDEQ(id)).Exist(ctx)
}

func (repo *academyRepository) Update(ctx context.Context, d *ent.Academy) error {
	clause := repo.db.Academy.Update().
		Where(
			academy.IDEQ(d.ID),
			academy.UserIDEQ(d.UserID),
		)

	mu := clause.Mutation()

	if d.AddressDetail == nil {
		mu.ClearAddressDetail()
	}

	clause.
		SetSigunguID(d.SigunguID).
		SetName(d.Name).
		SetCallNumber(d.CallNumber).
		SetAddressRoad(d.AddressRoad).
		SetNillableAddressDetail(d.AddressDetail).
		ClearYoga()

	if len(d.Edges.Yoga) > 0 {
		clause.AddYoga(d.Edges.Yoga...)
	}

	return clause.Exec(ctx)
}

func (repo *academyRepository) Patch(ctx context.Context, d *request.AcademyPatchBody, id, userId int) error {
	info := structs.New(d.Info)
	updateableData := utils.GetUpdateableDataV2(info, academy.Columns)
	delete(updateableData, "id")

	clause := repo.db.Academy.Update().
		Where(
			academy.IDEQ(id),
			academy.UserIDEQ(userId),
		)

	for key, val := range updateableData {
		clause.Mutation().SetField(key, val)
	}

	if d.YogaIDs != nil {
		yogaIds := *d.YogaIDs
		clause.ClearYoga().AddYogaIDs(yogaIds...)
	}

	return clause.Exec(ctx)
}

func (repo *academyRepository) Get(ctx context.Context, id int) (*ent.Academy, error) {
	var selectColumns []string
	for _, v := range academy.Columns {
		if v == academy.FieldBusinessCode {
			continue
		}
		selectColumns = append(selectColumns, v)
	}

	return repo.db.Academy.
		Query().
		Select(selectColumns...).
		WithAreaSigungu(
			func(asgq *ent.AreaSiGunguQuery) {
				asgq.Select(areasigungu.FieldName)
			},
		).
		WithYoga(
			func(yq *ent.YogaQuery) {
				yq.Select(yoga.FieldLevel, yoga.FieldNameKor)
			},
		).
		WithYogaRaw(
			func(yrq *ent.YogaRawQuery) {
				yrq.Select(yogaraw.FieldID, yogaraw.FieldName)
			},
		).
		Where(academy.IDEQ(id)).
		Only(ctx)
}

func (repo *academyRepository) GetOnlyIdByUserId(ctx context.Context, userId int) (id int, err error) {
	return repo.db.Academy.Query().Where(academy.UserIDEQ(userId)).OnlyID(ctx)
}

func (repo *academyRepository) Total(ctx context.Context, yogaIDs *[]int, sigunguID *int, academyName *string) (result int, err error) {
	clause := repo.db.Academy.Query()

	clause = repo.conditionQuery(ctx, yogaIDs, sigunguID, academyName, clause)
	result, err = clause.Count(ctx)
	return
}

func (repo *academyRepository) List(ctx context.Context,
	pgModule *utils.Pagination,
	yogaIDs *[]int, sigunguID *int, academyName *string,
	orderCol *string, orderType common.Sorts,
) (result []*ent.Academy, err error) {
	clause := repo.db.Academy.
		Query().
		WithAreaSigungu(
			func(asgq *ent.AreaSiGunguQuery) {
				asgq.Select(areasigungu.FieldName)
			},
		).
		WithYoga(
			func(yq *ent.YogaQuery) {
				yq.Select(yoga.FieldLevel, yoga.FieldNameKor)
			},
		).
		Limit(pgModule.GetLimit()).
		Offset(pgModule.GetOffset())

	useableOrderCol := map[string]string{
		"NAME": academy.FieldName,
		"ID":   academy.FieldID,
	}

	useableOrderFunc := map[common.Sorts]func(v ...string) ent.OrderFunc{
		common.DESC: ent.Desc,
		common.ASC:  ent.Asc,
	}

	if orderCol != nil {
		clause.Order(useableOrderFunc[orderType](useableOrderCol[*orderCol]))
	} else {
		clause.Order(useableOrderFunc[orderType](academy.FieldID))
	}

	clause = repo.conditionQuery(ctx, yogaIDs, sigunguID, academyName, clause)

	result, err = clause.All(ctx)
	return
}

func (repo *academyRepository) conditionQuery(
	ctx context.Context,
	yogaIDs *[]int,
	sigunguID *int,
	academyName *string,
	clause *ent.AcademyQuery,
) *ent.AcademyQuery {
	if academyName != nil {
		clause.Where(academy.NameContains(*academyName))
	}

	if sigunguID != nil {
		clause.Where(academy.HasAreaSigunguWith(areasigungu.IDEQ(*sigunguID)))
	}

	if yogaIDs != nil {
		clause.Where(academy.HasYogaWith(yoga.IDIn(*yogaIDs...)))
	}

	return clause
}
