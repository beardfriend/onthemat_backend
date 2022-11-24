package repository

import (
	"context"
	"errors"

	"onthemat/internal/app/common"
	"onthemat/internal/app/utils"
	"onthemat/pkg/ent"
	"onthemat/pkg/ent/predicate"
	"onthemat/pkg/ent/yoga"
	"onthemat/pkg/ent/yogagroup"
)

type YogaRepository interface {
	CreateGroup(ctx context.Context, data *ent.YogaGroup) error
	UpdateGroup(ctx context.Context, data *ent.YogaGroup) error
	DeleteGroups(ctx context.Context, ids []int) (int, error)
	GroupTotal(ctx context.Context, p *common.TotalParams) (count int, err error)
	GroupList(ctx context.Context, pgModule *utils.Pagination, p *common.ListParams) (result []*ent.YogaGroup, err error)
	Create(ctx context.Context, data *ent.Yoga) error
	Update(ctx context.Context, data *ent.Yoga) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, groupdId int) ([]*ent.Yoga, error)
}

type yogaRepository struct {
	db *ent.Client
}

func NewYogaRepository(db *ent.Client) YogaRepository {
	return &yogaRepository{
		db: db,
	}
}

// ------------------- Group -------------------

func (repo *yogaRepository) CreateGroup(ctx context.Context, data *ent.YogaGroup) error {
	return repo.db.YogaGroup.Create().
		SetCategory(data.Category).
		SetCategoryEng(data.CategoryEng).
		SetNillableDescription(&data.Description).
		Exec(ctx)
}

func (repo *yogaRepository) UpdateGroup(ctx context.Context, data *ent.YogaGroup) error {
	return repo.db.YogaGroup.UpdateOneID(data.ID).
		SetCategory(data.Category).
		SetCategoryEng(data.CategoryEng).
		SetNillableDescription(&data.Description).
		Exec(ctx)
}

func (repo *yogaRepository) GroupTotal(ctx context.Context, p *common.TotalParams) (count int, err error) {
	clause := repo.db.YogaGroup.Query()
	clause, err = repo.groupConditionQuery(p, clause)
	if err != nil {
		return
	}
	count, err = clause.Count(ctx)
	return
}

func (repo *yogaRepository) GroupList(ctx context.Context, pgModule *utils.Pagination, p *common.ListParams) (result []*ent.YogaGroup, err error) {
	clause := repo.db.YogaGroup.Query().
		Limit(pgModule.GetLimit()).
		Offset(pgModule.GetOffset())

	clause, err = repo.groupConditionQuery(&common.TotalParams{
		SearchKey:   p.SearchKey,
		SearchValue: p.SearchValue,
	}, clause)
	if err != nil {
		return
	}

	result, err = clause.All(ctx)
	return
}

func (repo *yogaRepository) groupConditionQuery(p *common.TotalParams, clause *ent.YogaGroupQuery) (*ent.YogaGroupQuery, error) {
	useableSearchCol := map[string]func(v string) predicate.YogaGroup{
		"NAME": yogagroup.CategoryContains,
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

func (repo *yogaRepository) DeleteGroups(ctx context.Context, ids []int) (int, error) {
	return repo.db.YogaGroup.Delete().Where(yogagroup.IDIn(ids...)).Exec(ctx)
}

// ------------------- Yoga -------------------

func (repo *yogaRepository) Create(ctx context.Context, data *ent.Yoga) error {
	return repo.db.Yoga.Create().
		SetNameKor(data.NameKor).
		SetNillableNameEng(data.NameEng).
		SetNillableLevel(data.Level).
		SetNillableDescription(data.Description).
		SetYogaGroupID(data.Edges.YogaGroup.ID).
		Exec(ctx)
}

func (repo *yogaRepository) Update(ctx context.Context, data *ent.Yoga) error {
	// 업데이트 쿼리는 전부 WHEre로바꾸기 아니면 select까지 해버림
	return repo.db.Yoga.Update().Where(yoga.IDEQ(data.ID)).
		SetNameKor(data.NameKor).
		SetNillableNameEng(data.NameEng).
		SetNillableDescription(data.Description).
		SetNillableLevel(data.Level).
		SetYogaGroupID(data.Edges.YogaGroup.ID).
		Exec(ctx)
}

func (repo *yogaRepository) Delete(ctx context.Context, id int) error {
	return repo.db.Yoga.DeleteOneID(id).Exec(ctx)
}

func (repo *yogaRepository) List(ctx context.Context, groupdId int) ([]*ent.Yoga, error) {
	// join reference
	return repo.db.YogaGroup.Query().
		Where(yogagroup.IDEQ(groupdId)).
		QueryYoga().
		Order(ent.Desc(yoga.FieldID)).
		All(ctx)
}
