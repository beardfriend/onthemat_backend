package repository

import (
	"context"
	"fmt"

	"onthemat/pkg/ent"
)

type AcademyRepository interface {
	Create(ctx context.Context) error
}

type academyRepository struct {
	db *ent.Client
}

func NewAcademyRepository(db *ent.Client) AcademyRepository {
	return &academyRepository{
		db: db,
	}
}

func (svc *academyRepository) Create(ctx context.Context) error {
	u, err := svc.db.User.Create().SetPhoneNum("010").SetNickname("nick").Save(ctx)
	if err != nil {
		return err
	}

	_, err = svc.db.UserNormal.Create().SetEmail("asd").SetPassword("pass").SetUser(u).Save(ctx)
	if err != nil {
		return err
	}
	_, err = svc.db.Acadmey.Create().SetName("test").SetUser(u).Save(ctx)
	return err
}

func (svc *academyRepository) CreateTransaction(ctx context.Context) error {
	tx, err := svc.db.Tx(ctx)
	if err != nil {
		return err
	}
	txClient := tx.Client()
	svc.db = txClient
	// Use the "Gen" below, but give it the transactional client; no code changes to "Gen".
	if err := svc.Create(ctx); err != nil {
		return rollback(tx, err)
	}
	return tx.Commit()
}

func rollback(tx *ent.Tx, err error) error {
	if rerr := tx.Rollback(); rerr != nil {
		err = fmt.Errorf("%w: %v", err, rerr)
	}
	return err
}
