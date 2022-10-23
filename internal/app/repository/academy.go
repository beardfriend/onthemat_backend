package repository

import (
	"context"

	"onthemat/pkg/ent"
	"onthemat/pkg/entx"
)

type AcademyRepository interface {
	Create(ctx context.Context, academy *ent.Acadmey) error
}

type academyRepository struct {
	db *ent.Client
}

func NewAcademyRepository(db *ent.Client) AcademyRepository {
	return &academyRepository{
		db: db,
	}
}

func (svc *academyRepository) Create(ctx context.Context, academy *ent.Acadmey) error {
	err := entx.WithTx(ctx, svc.db, func(tx *ent.Tx) error {
		user := academy.Edges.User

		u, err := tx.User.Create().
			SetEmail(user.Email).
			SetPassword(user.Password).
			SetNickname(user.Nickname).
			Save(ctx)
		if err != nil {
			return err
		}

		_, err = tx.Acadmey.Create().
			SetName(academy.Name).
			SetBusinessCode(*academy.BusinessCode).
			SetFullAddress(academy.FullAddress).
			SetUser(u).
			Save(ctx)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
