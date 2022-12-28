package repository

import (
	"context"

	"onthemat/pkg/ent"
	"onthemat/pkg/ent/areasigungu"
	"onthemat/pkg/entx"
)

type AreaRepository interface {
	Create(ctx context.Context, d *ent.AreaSiDo, sigungu []*ent.AreaSiGungu) error
	GetSigunGu(ctx context.Context, name string) (*ent.AreaSiGungu, error)
	GetSigunguIdByAdmCode(ctx context.Context, admCode string) (id int, err error)
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
			SetVersion(d.Version).AddSigungu(si...).
			Exec(ctx)
	})
}

func (repo *areaRepository) GetSigunGu(ctx context.Context, name string) (*ent.AreaSiGungu, error) {
	return repo.db.AreaSiGungu.Query().
		Where(areasigungu.NameEQ(name)).
		Only(ctx)
}

func (repo *areaRepository) GetSigunguIdByAdmCode(ctx context.Context, admCode string) (id int, err error) {
	return repo.db.AreaSiGungu.Query().Where(areasigungu.AdmCodeEQ(admCode)).OnlyID(ctx)
}
