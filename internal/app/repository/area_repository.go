package repository

import (
	"context"

	"onthemat/pkg/ent"
	"onthemat/pkg/entx"
)

type AreaRepository interface {
	Create(ctx context.Context, d *ent.AreaSiDo, sigungu []*ent.AreaSiGungu) error
}

type areaRepository struct {
	db *ent.Client
}

func NewAreaRepository(db *ent.Client) AreaRepository {
	return &areaRepository{
		db: db,
	}
}

func (repo *areaRepository) Create(ctx context.Context, d *ent.AreaSiDo, sigungu []*ent.AreaSiGungu) error {
	return entx.WithTx(ctx, repo.db, func(tx *ent.Tx) (err error) {
		bulk := make([]*ent.AreaSiGunguCreate, len(sigungu))

		for i, v := range sigungu {
			bulk[i] = repo.db.AreaSiGungu.Create().
				SetAdmCode(v.AdmCode).
				SetName(v.Name).
				SetVersion(d.Version).
				SetNillableParentCode(v.ParentCode)
		}
		si, err := repo.db.AreaSiGungu.CreateBulk(bulk...).Save(ctx)
		if err != nil {
			return
		}

		return repo.db.AreaSiDo.Create().
			SetName(d.Name).
			SetAdmCode(d.AdmCode).
			SetVersion(d.Version).AddSigungus(si...).
			Exec(ctx)
	})
}
