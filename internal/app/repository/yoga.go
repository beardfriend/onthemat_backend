package repository

import (
	"context"
	"errors"

	"onthemat/internal/app/common"
	"onthemat/internal/app/transport/request"

	"onthemat/internal/app/utils"
	"onthemat/pkg/ent"
	"onthemat/pkg/ent/predicate"
	"onthemat/pkg/ent/yoga"
	"onthemat/pkg/ent/yogagroup"
)

type YogaRepository interface {
	// test_success
	CreateGroup(ctx context.Context, data *ent.YogaGroup) error
	UpdateGroup(ctx context.Context, data *ent.YogaGroup) error
	PatchGroup(ctx context.Context, data *request.YogaGroupPatchBody, id int) error
	DeleteGroups(ctx context.Context, ids []int) (int, error)

	Create(ctx context.Context, data *ent.Yoga) error
	Update(ctx context.Context, data *ent.Yoga) error
	Patch(ctx context.Context, data *request.YogaPatchBody, id int) error
	Delete(ctx context.Context, id int) error

	// not yet
	GroupTotal(ctx context.Context, p *common.TotalParams) (count int, err error)
	GroupList(ctx context.Context, pgModule *utils.Pagination, p *common.ListParams) (result []*ent.YogaGroup, err error)

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
		SetNillableDescription(data.Description).
		Exec(ctx)
}

func (repo *yogaRepository) UpdateGroup(ctx context.Context, data *ent.YogaGroup) error {
	clause := repo.db.YogaGroup.Update().Where(yogagroup.IDEQ(data.ID))
	mu := clause.Mutation()

	if data.Description == nil {
		mu.ClearDescription()
	}

	return clause.
		SetCategory(data.Category).
		SetCategoryEng(data.CategoryEng).
		SetNillableDescription(data.Description).
		Exec(ctx)
}

func (repo *yogaRepository) PatchGroup(ctx context.Context, data *request.YogaGroupPatchBody, id int) error {
	updateableData := utils.GetUpdateableData(data, yogagroup.Columns)

	clause := repo.db.YogaGroup.Update().Where(yogagroup.IDEQ(id))
	for key, val := range updateableData {
		clause.Mutation().SetField(key, val)
	}
	return clause.Exec(ctx)
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
		SetYogaGroupID(data.YogaGroupID).
		Exec(ctx)
}

func (repo *yogaRepository) Update(ctx context.Context, data *ent.Yoga) error {
	clause := repo.db.Yoga.Update().Where(yoga.IDEQ(data.ID))
	mu := clause.Mutation()

	if data.Level == nil {
		mu.ClearLevel()
	}
	if data.NameEng == nil {
		mu.ClearNameEng()
	}
	if data.Description == nil {
		mu.ClearDescription()
	}
	return clause.
		SetNillableDescription(data.Description).
		SetNillableLevel(data.Level).
		SetNillableNameEng(data.NameEng).
		SetNameKor(data.NameKor).
		SetYogaGroupID(data.YogaGroupID).Exec(ctx)
}

func (repo *yogaRepository) Patch(ctx context.Context, data *request.YogaPatchBody, id int) error {
	updateableData := utils.GetUpdateableData(data, yoga.Columns)

	clause := repo.db.Yoga.Update().Where(yoga.IDEQ(id))
	for key, val := range updateableData {
		clause.Mutation().SetField(key, val)
	}
	return clause.Exec(ctx)
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
