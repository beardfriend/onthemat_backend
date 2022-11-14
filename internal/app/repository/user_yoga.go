package repository

import (
	"context"

	"onthemat/pkg/ent"
	"onthemat/pkg/ent/useryoga"
)

type UserYogaRepository interface {
	CreateMany(ctx context.Context, value []*ent.UserYoga, userId int) ([]*ent.UserYoga, error)
	DeleteMany(ctx context.Context, Ids []int) (int, error)
}

type userYogaRepository struct {
	db *ent.Client
}

func NewUserYogaRepository(db *ent.Client) UserYogaRepository {
	return &userYogaRepository{
		db: db,
	}
}

func (repo *userYogaRepository) CreateMany(ctx context.Context, value []*ent.UserYoga, userId int) ([]*ent.UserYoga, error) {
	bulk := make([]*ent.UserYogaCreate, len(value))
	for i, v := range value {
		bulk[i] = repo.db.UserYoga.Create().SetName(v.Name).SetUserType(v.UserType).SetUserID(userId)
	}
	return repo.db.UserYoga.CreateBulk(bulk...).Save(ctx)
}

func (repo *userYogaRepository) DeleteMany(ctx context.Context, Ids []int) (int, error) {
	return repo.db.UserYoga.Delete().Where(
		useryoga.IDIn(Ids...),
	).Exec(ctx)
}
