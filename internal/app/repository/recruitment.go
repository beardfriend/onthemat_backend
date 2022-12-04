package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"onthemat/internal/app/transport"
	"onthemat/internal/app/transport/request"
	"onthemat/internal/app/utils"
	"onthemat/pkg/ent"
	"onthemat/pkg/ent/academy"
	"onthemat/pkg/ent/areasigungu"
	"onthemat/pkg/ent/recruitment"
	ri "onthemat/pkg/ent/recruitmentinstead"
	"onthemat/pkg/ent/teacher"
	"onthemat/pkg/ent/yoga"
	"onthemat/pkg/entx"

	"entgo.io/ent/dialect/sql"
	"github.com/fatih/structs"
)

type RecruitmentRepository interface {
	Create(ctx context.Context, d *ent.Recruitment) (err error)
	Update(ctx context.Context, d *ent.Recruitment) (err error)
	Patch(ctx context.Context, d *request.RecruitmentPatchBody, id, academyId int) (isCreated bool, err error)
	PatchDeletedAt(ctx context.Context, id, academyId int) (err error)
	Total(ctx context.Context, startDateTime, endDateTime *transport.TimeString, yogaIds, sigunguId *[]int) (result int, err error)
	List(ctx context.Context, pgModule *utils.Pagination, startDateTime, endDateTime *transport.TimeString, yogaIds, sigunguId *[]int) ([]*ent.Recruitment, error)
	Exist(ctx context.Context, id int) (bool, error)
	Get(ctx context.Context, id int) (*ent.Recruitment, error)
}

type recruitmentRepository struct {
	db                 *ent.Client
	recruitInsteadRepo *recruitmentInsteadRepo
}

func NewRecruitmentRepository(db *ent.Client) RecruitmentRepository {
	return &recruitmentRepository{
		db: db,
	}
}

func (repo *recruitmentRepository) Create(ctx context.Context, d *ent.Recruitment) (err error) {
	return entx.WithTx(ctx, repo.db, func(tx *ent.Tx) (err error) {
		client := tx.Client()

		recruit, err := client.Recruitment.Create().
			SetWriterID(d.AcademyID).
			Save(ctx)
		if err != nil {
			return
		}

		if len(d.Edges.RecruitmentInstead) > 0 {
			err = repo.recruitInsteadRepo.createMany(ctx, client, d.Edges.RecruitmentInstead, recruit.ID)
			if err != nil {
				return
			}
		}
		return
	})
}

func (repo *recruitmentRepository) PatchDeletedAt(ctx context.Context, id, academyId int) (err error) {
	rowAffcetd, err := repo.db.Recruitment.Update().
		SetDeletedAt(transport.TimeString(time.Now())).
		Where(recruitment.AcademyIDEQ(academyId), recruitment.IDEQ(id)).
		Save(ctx)
	if err != nil {
		return
	}

	if rowAffcetd != 1 {
		err = errors.New(ErrOnlyOwnUser)
		return
	}
	return
}

const ErrOnlyOwnUser = "소유자만 접근할 수 있습니다"

func (repo *recruitmentRepository) Update(ctx context.Context, d *ent.Recruitment) (err error) {
	return entx.WithTx(ctx, repo.db.Debug(), func(tx *ent.Tx) (err error) {
		client := tx.Client()

		rowAffected, err := client.Recruitment.Update().
			Where(
				recruitment.IDEQ(d.ID),
				recruitment.AcademyIDEQ(d.AcademyID),
			).
			// 필요할까 ?
			SetWriterID(d.AcademyID).
			SetIsFinish(d.IsFinish).
			SetIsOpen(d.IsOpen).
			Save(ctx)
		if err != nil {
			return
		}

		if rowAffected != 1 {
			err = errors.New(ErrOnlyOwnUser)
			return
		}

		if len(d.Edges.RecruitmentInstead) > 0 {
			ids, err := repo.recruitInsteadRepo.getIdsByRecruitId(ctx, client, d.ID)
			if err != nil {
				return err
			}

			requestIds := repo.extractIdsFromInsteadrepo(d.Edges.RecruitmentInstead)
			createable, updateable, deleteable := utils.MakeDataForCondition(requestIds, ids)

			if len(createable) > 0 {
				createData := repo.filterRecruitInstead(d.Edges.RecruitmentInstead, createable)
				err = repo.recruitInsteadRepo.createMany(ctx, client, createData, d.ID)
				if err != nil {
					return err
				}
			}

			if len(updateable) > 0 {
				updateData := repo.filterRecruitInstead(d.Edges.RecruitmentInstead, updateable)
				err = repo.recruitInsteadRepo.updateMany(ctx, client, updateData, d.ID)
				if err != nil {
					return err
				}
			}

			if len(deleteable) > 0 {
				err = repo.recruitInsteadRepo.deleteByIds(ctx, client, deleteable)
				if err != nil {
					return err
				}
			}

		} else {

			err = repo.recruitInsteadRepo.deleteByRecruitId(ctx, client, d.ID)
			if err != nil {
				return err
			}

		}
		return
	})
}

func (repo *recruitmentRepository) Patch(ctx context.Context, d *request.RecruitmentPatchBody, id, academyId int) (isCreated bool, err error) {
	entx.WithTx(ctx, repo.db, func(tx *ent.Tx) (err error) {
		client := tx.Client()

		clause := client.Recruitment.Update().
			Where(recruitment.IDEQ(id), recruitment.AcademyID(academyId))

		if d.Info != nil {
			recruitInfo := structs.New(d.Info)
			updateableRecruitInfo := utils.GetUpdateableDataV2(recruitInfo, recruitment.Columns)
			for key, val := range updateableRecruitInfo {
				clause.Mutation().SetField(key, val)
			}
		}
		rowAffected, err := clause.Save(ctx)
		if err != nil {
			return
		}

		if rowAffected != 1 {
			err = errors.New(ErrOnlyOwnUser)
			return
		}

		if d.InsteadInfo != nil {
			for _, v := range d.InsteadInfo {
				s := structs.New(v)
				res := utils.GetUpdateableDataV2(s, ri.Columns)

				c := client.RecruitmentInstead

				// Update
				if v.ID != nil {
					u := c.Update().Where(ri.IDEQ(*v.ID))

					for key, val := range res {
						u.Mutation().SetField(key, val)
					}

					if err = u.Exec(ctx); err != nil {
						return
					}
					// Create
				} else {
					cr := c.Create()

					for key, val := range res {
						cr.Mutation().SetField(key, val)
					}
					if err = cr.Exec(ctx); err != nil {
						return
					}
					isCreated = true
				}
			}
		}
		return
	})
	return
}

func (repo *recruitmentRepository) Get(ctx context.Context, id int) (*ent.Recruitment, error) {
	return repo.db.Debug().Recruitment.Query().
		WithRecruitmentInstead(
			func(riq *ent.RecruitmentInsteadQuery) {
				riq.WithYoga(
					func(yq *ent.YogaQuery) {
						yq.Select(yoga.FieldID, yoga.FieldNameKor)
					},
				)
				riq.WithApplicant(
					func(tq *ent.TeacherQuery) {
						tq.Select(teacher.FieldID)
					})
			}).
		WithWriter().
		Where(recruitment.DeletedAtIsNil(), recruitment.IDEQ(id)).
		Only(ctx)
}

func (repo *recruitmentRepository) Total(
	ctx context.Context,
	startDateTime,
	endDateTime *transport.TimeString,
	yogaIds,
	sigunguId *[]int,
) (result int, err error) {
	clause := repo.db.Recruitment.Query()

	clause = repo.conditionQuery(clause, startDateTime, endDateTime, yogaIds, sigunguId)
	result, err = clause.Count(ctx)
	return
}

func (repo *recruitmentRepository) List(
	ctx context.Context,
	pgModule *utils.Pagination,
	startDateTime,
	endDateTime *transport.TimeString,
	yogaIds,
	sigunguId *[]int,
) ([]*ent.Recruitment, error) {
	clause := repo.db.Debug().Recruitment.Query().
		WithRecruitmentInstead(
			func(riq *ent.RecruitmentInsteadQuery) {
				riq.Select(ri.FieldID, ri.FieldRecruitmentID, ri.FieldSchedule)
				riq.WithYoga(
					func(yq *ent.YogaQuery) {
						yq.Select(yoga.FieldLevel, yoga.FieldNameKor)
					},
				)
			},
		).
		WithWriter(func(aq *ent.AcademyQuery) {
			aq.Select(academy.FieldName, academy.FieldSigunguID)
			aq.WithAreaSigungu(func(asgq *ent.AreaSiGunguQuery) {
				asgq.Select(areasigungu.FieldName)
			})
		}).
		Where(recruitment.DeletedAtIsNil()).
		Limit(pgModule.GetLimit()).
		Offset(pgModule.GetOffset()).
		Order(ent.Desc(recruitment.FieldUpdatedAt))

	clause = repo.conditionQuery(clause, startDateTime, endDateTime, yogaIds, sigunguId)
	return clause.All(ctx)
}

func (repo *recruitmentRepository) Exist(ctx context.Context, id int) (bool, error) {
	return repo.db.Recruitment.Query().Where(recruitment.IDEQ(id)).Exist(ctx)
}

// 시간으로 조회 // 요가로 조회 // 학원 위치로 조회
func (repo *recruitmentRepository) conditionQuery(
	clause *ent.RecruitmentQuery,
	startDateTime, endDateTime *transport.TimeString,
	yogaIds, sigunguId *[]int,
) *ent.RecruitmentQuery {
	if startDateTime != nil && endDateTime != nil {
		clause.Where(
			recruitment.HasRecruitmentInsteadWith(
				func(s *sql.Selector) {
					tableName := fmt.Sprintf(`%s, jsonb_array_elements(%s) c`, ri.Table, ri.FieldSchedule)
					s.From(sql.Table(tableName))

					p := &sql.Predicate{}
					p.Append(func(b *sql.Builder) {
						// startDateTime으로 조회
						b.WriteString("c ->> 'startDateTime'").
							WriteOp(sql.OpGTE).Arg(startDateTime.ToString())

						b.WriteString(" AND ")

						// endDateTime으로 조회
						b.WriteString("c ->>'startDateTime'").
							WriteOp(sql.OpLTE).
							Arg(endDateTime.ToString())
					})
					s.Where(p)
				},
			),
		)
	}
	if yogaIds != nil {
		clause.Where(
			recruitment.HasRecruitmentInsteadWith(
				ri.HasYogaWith(yoga.IDIn(*yogaIds...)),
			),
		)
	}

	if sigunguId != nil {
		clause.Where(
			recruitment.HasWriterWith(
				academy.SigunguIDIn(*sigunguId...),
			),
		)
	}
	return clause
}

func (repo *recruitmentRepository) extractIdsFromInsteadrepo(vals []*ent.RecruitmentInstead) []int {
	var result []int
	for _, s := range vals {
		result = append(result, s.ID)
	}
	return result
}

func (repo *recruitmentRepository) filterRecruitInstead(vals []*ent.RecruitmentInstead, ids []int) []*ent.RecruitmentInstead {
	result := make([]*ent.RecruitmentInstead, 0)

	for k := 0; k < len(ids); k++ {
		for i := k; i < len(vals); i++ {
			if vals[i].ID == ids[k] {
				result = append(result, vals[i])
				break
			}
		}
	}

	return result
}
