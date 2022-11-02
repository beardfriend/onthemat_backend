package repository

import (
	"context"

	"onthemat/pkg/ent"
	"onthemat/pkg/ent/user"
)

type UserRepository interface {
	Create(ctx context.Context, user *ent.User) (*ent.User, error)
	Update(ctx context.Context, user *ent.User) (*ent.User, error)
	GetBySocialKey(ctx context.Context, u *ent.User) (*ent.User, error)
	Get(ctx context.Context, id int) (*ent.User, error)
}

type userRepository struct {
	db *ent.Client
}

func NewUserRepository(db *ent.Client) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (repo *userRepository) Get(ctx context.Context, id int) (*ent.User, error) {
	return repo.db.User.Get(ctx, id)
}

func (repo *userRepository) Create(ctx context.Context, user *ent.User) (*ent.User, error) {
	u, err := repo.db.User.Create().
		SetNillableEmail(&user.Email).
		SetNillablePassword(&user.Password).
		SetNillableSocialKey(&user.SocialKey).
		SetNillableSocialName(&user.SocialName).
		SetNillableTermAgreeAt(&user.TermAgreeAt).
		SetNillablePhoneNum(&user.PhoneNum).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (repo *userRepository) Update(ctx context.Context, user *ent.User) (*ent.User, error) {
	u, err := repo.db.User.UpdateOneID(user.ID).
		SetNillableEmail(&user.Email).
		SetNillableTermAgreeAt(&user.TermAgreeAt).
		SetNillableType(user.Type).
		SetNillablePhoneNum(&user.PhoneNum).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (repo *userRepository) GetBySocialKey(ctx context.Context, u *ent.User) (*ent.User, error) {
	return repo.db.Debug().User.
		Query().
		Where(
			user.SocialKeyEQ(u.SocialKey),
			user.SocialNameEQ(u.SocialName)).Only(ctx)
}
