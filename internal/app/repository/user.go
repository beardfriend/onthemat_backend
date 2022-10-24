package repository

import (
	"context"

	"onthemat/pkg/ent"
)

type UserRepository interface {
	Create(ctx context.Context, user *ent.User) error
}

type userRepository struct {
	db *ent.Client
}

func NewUserRepository(db *ent.Client) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (repo *userRepository) Create(ctx context.Context, user *ent.User) error {
	if _, err := repo.db.User.Create().
		SetEmail(user.Email).
		SetPassword(user.Password).
		Save(ctx); err != nil {
		return err
	}
	return nil
}
