package repository

import (
	"context"

	"onthemat/pkg/ent"
)

type UserYogaRepository interface {
	CreateMany(ctx context.Context, value []*ent.UserYoga, userId int) error
}

type userYogaRepository struct {
	db *ent.Client
}

func NewUserYogaRepository(db *ent.Client) UserYogaRepository {
	return &userYogaRepository{
		db: db,
	}
}

func (repo *userYogaRepository) CreateMany(ctx context.Context, value []*ent.UserYoga, userId int) error {
	bulk := make([]*ent.UserYogaCreate, len(value))
	for i, v := range value {
		bulk[i] = repo.db.UserYoga.Create().SetName(v.Name).SetUserType(v.UserType).SetUserID(userId)
	}
	return repo.db.UserYoga.CreateBulk(bulk...).Exec(ctx)
}
