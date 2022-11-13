package repository

import (
	"context"

	"onthemat/pkg/ent"
	"onthemat/pkg/ent/user"
)

type UserRepository interface {
	Create(ctx context.Context, user *ent.User) (*ent.User, error)
	Update(ctx context.Context, user *ent.User) (*ent.User, error)
	// UpdatePassword()
	UpdateTempPassword(ctx context.Context, u *ent.User) error
	UpdateEmailVerifeid(ctx context.Context, userId int) error
	GetBySocialKey(ctx context.Context, u *ent.User) (*ent.User, error)
	GetByEmail(ctx context.Context, email string) (*ent.User, error)
	GetByEmailPassword(ctx context.Context, u *ent.User) (*ent.User, error)
	FindByEmail(ctx context.Context, email string) (bool, error)
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
	return repo.db.User.Query().Where(user.IDEQ(id)).Only(ctx)
}

func (repo *userRepository) Create(ctx context.Context, user *ent.User) (*ent.User, error) {
	u, err := repo.db.User.Create().
		SetNillableNickname(user.Nickname).
		SetNillableEmail(user.Email).
		SetNillablePassword(user.Password).
		SetNillableSocialKey(user.SocialKey).
		SetNillableSocialName(user.SocialName).
		SetNillableTermAgreeAt(user.TermAgreeAt).
		SetNillablePhoneNum(user.PhoneNum).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (repo *userRepository) Update(ctx context.Context, user *ent.User) (*ent.User, error) {
	u, err := repo.db.User.UpdateOneID(user.ID).
		SetNillableEmail(user.Email).
		SetNillableNickname(user.Nickname).
		SetNillablePhoneNum(user.PhoneNum).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (repo *userRepository) UpdateTempPassword(ctx context.Context, u *ent.User) error {
	return repo.db.User.Update().
		SetNillableTempPassword(u.TempPassword).Where(
		user.EmailEQ(*u.Email),
	).Exec(ctx)
}

func (repo *userRepository) UpdateEmailVerifeid(ctx context.Context, userId int) error {
	return repo.db.User.UpdateOneID(userId).
		SetIsEmailVerified(true).Exec(ctx)
}

func (repo *userRepository) GetBySocialKey(ctx context.Context, u *ent.User) (*ent.User, error) {
	return repo.db.User.
		Query().
		Where(
			user.SocialKeyEQ(*u.SocialKey),
			user.SocialNameEQ(*u.SocialName)).Only(ctx)
}

func (repo *userRepository) GetByEmail(ctx context.Context, email string) (*ent.User, error) {
	return repo.db.User.
		Query().
		Where(
			user.EmailEQ(email)).Only(ctx)
}

func (repo *userRepository) GetByEmailPassword(ctx context.Context, u *ent.User) (*ent.User, error) {
	return repo.db.User.
		Query().
		Where(
			user.Or(user.PasswordEQ(*u.Password), user.TempPasswordEQ(*u.Password)),
			user.EmailEQ(*u.Email),
		).Only(ctx)
}

func (repo *userRepository) FindByEmail(ctx context.Context, email string) (bool, error) {
	return repo.db.User.Query().Where(
		user.EmailEQ(email),
	).Exist(ctx)
}
