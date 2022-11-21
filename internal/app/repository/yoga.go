package repository

import (
	"context"

	"onthemat/pkg/ent"
	"onthemat/pkg/ent/yogagroup"
)

type YogaRepository interface {
	CreateGroup(ctx context.Context, data *ent.YogaGroup) error
	UpdateGroup(ctx context.Context, data *ent.YogaGroup) error
	DeleteGroups(ctx context.Context, ids []int) (int, error)
	Create(ctx context.Context, data *ent.Yoga) error
	CreateRaw(ctx context.Context, data *ent.YogaRaw) error
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

func (repo *yogaRepository) DeleteGroups(ctx context.Context, ids []int) (int, error) {
	return repo.db.YogaGroup.Delete().Where(yogagroup.IDIn(ids...)).Exec(ctx)
}

// ------------------- Yoga -------------------

func (repo *yogaRepository) Create(ctx context.Context, data *ent.Yoga) error {
	return repo.db.Yoga.Create().
		SetNameKor(data.NameKor).
		SetNillableNameEng(data.NameEng).
		SetNillableDescription(data.Description).
		SetYogaGroupID(data.Edges.YogaGroup.ID).
		Exec(ctx)
}

func (repo *yogaRepository) Update(ctx context.Context, data *ent.Yoga) error {
	return repo.db.Yoga.UpdateOneID(data.ID).
		SetNameKor(data.NameKor).
		SetNillableNameEng(data.NameEng).
		SetNillableDescription(data.Description).
		SetYogaGroupID(data.Edges.YogaGroup.ID).
		Exec(ctx)
}

func (repo *yogaRepository) UpdateYogaGroupID(ctx context.Context, id int, groupId int) error {
	return repo.db.Yoga.UpdateOneID(id).
		SetYogaGroupID(groupId).
		Exec(ctx)
}

func (repo *yogaRepository) Delete(ctx context.Context, id int) error {
	return repo.db.Yoga.DeleteOneID(id).Exec(ctx)
}

// ------------------- Yoga_raw -------------------

func (repo *yogaRepository) CreateRaw(ctx context.Context, data *ent.YogaRaw) error {
	return repo.db.YogaRaw.Create().
		SetName(data.Name).
		SetIsMigrated(false).
		SetUserID(data.Edges.User.ID).
		Exec(ctx)
}
