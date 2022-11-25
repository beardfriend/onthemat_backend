package repository

import (
	"context"
	"errors"
	"fmt"

	"onthemat/internal/app/model"
	"onthemat/internal/app/utils"

	"onthemat/pkg/ent"
	"onthemat/pkg/ent/academy"
	"onthemat/pkg/ent/areasigungu"
	"onthemat/pkg/ent/user"
	"onthemat/pkg/ent/yoga"
	"onthemat/pkg/entx"
)

type AcademyRepository interface {
	Create(ctx context.Context, d *ent.Academy) error
	Update(ctx context.Context, aca *ent.Academy, userId int, id int) error
	Get(ctx context.Context, userId int) (*ent.Academy, error)
	List(ctx context.Context,
		pgModule *utils.Pagination,
		yogaIDs *[]int, sigunguID *int, academyName *string,
		orderCol *string, orderType *string) (result []*ent.Academy, err error)

	Total(ctx context.Context, yogaIDs *[]int, sigunguID *int, academyName *string) (result int, err error)
	Patch(ctx context.Context, d *ent.Academy) error
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

func (repo *academyRepository) Create(ctx context.Context, d *ent.Academy) error {
	return entx.WithTx(ctx, repo.db, func(tx *ent.Tx) (err error) {
		err = repo.db.Academy.Create().
			SetNillableAddressDetail(d.AddressDetail).
			SetName(d.Name).
			SetBusinessCode(d.BusinessCode).
			SetCallNumber(d.CallNumber).
			SetAddressRoad(d.AddressRoad).
			SetUserID(d.UserID).
			SetSigunguID(d.SigunguID).
			Exec(ctx)
		if err != nil {
			return
		}

		if len(d.Edges.Yoga) > 0 {
			err = tx.Academy.Update().Where(
				academy.UserIDEQ(d.UserID),
			).AddYoga(d.Edges.Yoga...).Exec(ctx)

			if err != nil {
				return
			}
		}

		err = repo.db.User.Update().
			SetType(model.AcademyType).
			Where(user.IDEQ(d.UserID)).
			Exec(ctx)
		if err != nil {
			return
		}

		return
	})
}

func (repo *academyRepository) Patch(ctx context.Context, d *ent.Academy) error {
	dataForPatch := utils.StructToMap[ent.Value](d)

	delete(dataForPatch, "id")
	result := utils.MakeUseableFieldWithData(dataForPatch, academy.Columns)
	fmt.Println(result)
	clause := repo.db.Academy.Update()
	for key, val := range result {
		if val == nil {
			clause.Mutation().FieldCleared(key)
		}
		clause.Mutation().SetField(key, val)
	}

	if len(d.Edges.Yoga) > 0 {
		clause.ClearYoga().AddYoga(d.Edges.Yoga...)
	}

	return clause.Where(academy.IDEQ(d.ID)).Exec(ctx)
}

func (repo *academyRepository) Update(ctx context.Context, aca *ent.Academy, userId int, id int) error {
	return entx.WithTx(ctx, repo.db, func(tx *ent.Tx) (err error) {
		err = tx.Academy.UpdateOneID(userId).SetName(aca.Name).
			SetCallNumber(aca.CallNumber).
			SetAddressRoad(aca.AddressRoad).
			SetSigunguID(aca.SigunguID).
			Exec(ctx)
		if err != nil {
			return
		}

		err = tx.Academy.Update().Where(
			academy.IDEQ(aca.ID),
		).ClearYoga().Exec(ctx)
		if err != nil {
			return
		}

		tx.Academy.Update().Where(
			academy.IDEQ(aca.ID),
		).AddYoga(aca.Edges.Yoga...).Exec(ctx)
		if err != nil {
			return
		}

		return
	})
}

func (repo *academyRepository) Get(ctx context.Context, userId int) (*ent.Academy, error) {
	return repo.db.Academy.
		Query().
		Select(
			academy.FieldID,
			academy.FieldCallNumber,
			academy.FieldName,
			academy.FieldAddressRoad,
			academy.FieldAddressDetail,
			academy.FieldCreatedAt,
			academy.FieldUpdatedAt,
			academy.FieldSigunguID,
		).
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
		Where(academy.ID(userId)).
		Only(ctx)
}

func (repo *academyRepository) Total(ctx context.Context, yogaIDs *[]int, sigunguID *int, academyName *string) (result int, err error) {
	clause := repo.db.Academy.Query()

	clause = repo.conditionQuery(ctx, yogaIDs, sigunguID, academyName, clause)
	result, err = clause.Count(ctx)
	return
}

// TODO :페이지네이션 모듈 생성을 여기서 하지 않고 Usecase에서만 할 수 있도록
func (repo *academyRepository) List(ctx context.Context,
	pgModule *utils.Pagination,
	yogaIDs *[]int, sigunguID *int, academyName *string,
	orderCol *string, orderType *string,
) (result []*ent.Academy, err error) {
	clause := repo.db.Debug().Academy.
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
		"ID":        academy.FieldID,
		"CREATEDAT": academy.FieldCreatedAt,
	}

	if orderCol == nil || orderType == nil {
		clause.Order(ent.Desc(academy.FieldID))
	}

	if orderCol != nil && orderType != nil {
		orderCol := useableOrderCol[*orderCol]

		if orderCol == "" {
			err = errors.New(ErrOrderColumnInvalid)
			return
		}

		if *orderType == "ASC" {
			clause.Order(ent.Asc(orderCol))
		} else {
			clause.Order(ent.Desc(orderCol))
		}
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

func (repo *academyRepository) UseYoga(ctx context.Context, yogaIds []int, userID int) error {
	return repo.db.Academy.
		Update().
		Where(
			academy.UserIDEQ(userID),
		).AddYogaIDs(yogaIds...).Exec(ctx)
}

func (repo *academyRepository) UnUseYoga(ctx context.Context, yogaIds []int, userID int) error {
	return repo.db.Academy.
		Update().
		Where(
			academy.UserIDEQ(userID),
		).RemoveYogaIDs(yogaIds...).Exec(ctx)
}
