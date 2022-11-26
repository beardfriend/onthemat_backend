package repository

import (
	"context"

	"onthemat/internal/app/common"
	"onthemat/internal/app/transport/request"

	"onthemat/internal/app/utils"
	"onthemat/pkg/ent"
	"onthemat/pkg/ent/yoga"
	"onthemat/pkg/ent/yogagroup"
)

type YogaRepository interface {
	CreateGroup(ctx context.Context, data *ent.YogaGroup) error
	UpdateGroup(ctx context.Context, data *ent.YogaGroup) error
	PatchGroup(ctx context.Context, data *request.YogaGroupPatchBody, id int) error
	DeleteGroups(ctx context.Context, ids []int) (int, error)
	GroupTotal(ctx context.Context, category *string) (count int, err error)
	GroupList(ctx context.Context, pgModule *utils.Pagination, category *string, sorts common.Sorts) (result []*ent.YogaGroup, err error)

	Create(ctx context.Context, data *ent.Yoga) error
	Update(ctx context.Context, data *ent.Yoga) error
	Patch(ctx context.Context, data *request.YogaPatchBody, id int) error
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

func (repo *yogaRepository) GroupTotal(ctx context.Context, category *string) (count int, err error) {
	clause := repo.db.YogaGroup.Query()
	clause = repo.groupConditionQuery(category, clause)

	count, err = clause.Count(ctx)
	return
}

func (repo *yogaRepository) GroupList(ctx context.Context, pgModule *utils.Pagination, category *string, sorts common.Sorts) (result []*ent.YogaGroup, err error) {
	clause := repo.db.YogaGroup.Query().
		Limit(pgModule.GetLimit()).
		Offset(pgModule.GetOffset())

	if sorts == common.ASC {
		clause = clause.Order(ent.Asc(yogagroup.FieldID))
	} else {
		clause = clause.Order(ent.Desc(yogagroup.FieldID))
	}

	clause = repo.groupConditionQuery(category, clause)

	result, err = clause.All(ctx)
	return
}

func (repo *yogaRepository) groupConditionQuery(category *string, clause *ent.YogaGroupQuery) *ent.YogaGroupQuery {
	if category != nil {
		clause.Where(yogagroup.CategoryContains(*category))
	}

	return clause
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
