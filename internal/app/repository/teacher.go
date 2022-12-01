package repository

import (
	"context"

	"onthemat/internal/app/model"
	"onthemat/internal/app/transport/request"
	"onthemat/internal/app/utils"
	"onthemat/pkg/ent"
	"onthemat/pkg/ent/teacher"

	twe "onthemat/pkg/ent/teacherworkexperience"
	"onthemat/pkg/ent/user"

	"onthemat/pkg/entx"

	"github.com/fatih/structs"
)

type TeacherRepository interface {
	Create(ctx context.Context, d *ent.Teacher) error
	Update(ctx context.Context, d *ent.Teacher) (err error)
	Exist(ctx context.Context, id int) (bool, error)
	Patch(ctx context.Context, d *request.TeacherPatchBody, id, userId int) (err error)
	GetOnlyIdByUserId(ctx context.Context, userId int) (id int, err error)
}

type teacherRepository struct {
	db      *ent.Client
	workExp *teacherWorkExperience
	certifi *teacherCertification
	yoga    *yogaRepository
}

func NewTeacherRepository(db *ent.Client) TeacherRepository {
	workExp := new(teacherWorkExperience)
	certifi := new(teacherCertification)
	yoga := new(yogaRepository)
	return &teacherRepository{
		db:      db,
		workExp: workExp,
		certifi: certifi,
		yoga:    yoga,
	}
}

// success
func (repo *teacherRepository) Create(ctx context.Context, d *ent.Teacher) (err error) {
	return entx.WithTx(ctx, repo.db, func(tx *ent.Tx) (err error) {
		client := tx.Client()

		clause := client.Teacher.Create().
			SetUserID(d.UserID).
			SetNillableProfileImageUrl(d.ProfileImageUrl).
			SetName(d.Name).
			SetNillableAge(d.Age).
			SetNillableIntroduce(d.Introduce)

		if d.ID != 0 {
			clause.SetID(d.ID)
		}

		if len(d.Edges.Yoga) > 0 {
			clause.AddYoga(d.Edges.Yoga...)
		}

		if len(d.Edges.Sigungu) > 0 {
			clause.AddSigungu(d.Edges.Sigungu...)
		}

		teacher, err := clause.Save(ctx)
		if err != nil {
			return
		}
		if len(d.Edges.WorkExperience) > 0 {
			err = repo.workExp.createMany(ctx, client, d.Edges.WorkExperience, teacher.ID)
			if err != nil {
				return
			}
		}

		if len(d.Edges.Certification) > 0 {
			err = repo.certifi.createMany(ctx, client, d.Edges.Certification, teacher.ID)
			if err != nil {
				return
			}
		}

		if len(d.Edges.YogaRaw) > 0 {
			err = repo.yoga.createRaws(ctx, client, d.Edges.YogaRaw, &teacher.ID, nil)
			if err != nil {
				return
			}
		}

		return client.User.Update().
			SetType(model.TeacherType).
			Where(user.IDEQ(d.UserID)).
			Exec(ctx)
	})
}

func (repo *teacherRepository) Update(ctx context.Context, d *ent.Teacher) (err error) {
	return entx.WithTx(ctx, repo.db.Debug(), func(tx *ent.Tx) (err error) {
		client := tx.Client()
		clause := client.Teacher.Update().
			Where(
				teacher.IDEQ(d.ID),
				teacher.UserIDEQ(d.UserID)).
			SetUserID(d.UserID).
			SetNillableAge(d.Age).
			SetNillableIntroduce(d.Introduce).
			SetNillableProfileImageUrl(d.ProfileImageUrl).
			SetName(d.Name).
			SetIsProfileOpen(d.IsProfileOpen).
			ClearYoga().
			ClearSigungu()

		mu := clause.Mutation()

		if d.Age == nil {
			mu.ClearAge()
		}
		if d.Introduce == nil {
			mu.ClearIntroduce()
		}
		if d.ProfileImageUrl == nil {
			mu.ClearProfileImageUrl()
		}

		if len(d.Edges.Yoga) > 0 {
			clause.AddYoga(d.Edges.Yoga...)
		}

		if len(d.Edges.Sigungu) > 0 {
			clause.AddSigungu(d.Edges.Sigungu...)
		}

		err = clause.Exec(ctx)
		if err != nil {
			return
		}

		// WorkExperience
		if len(d.Edges.WorkExperience) > 0 {
			ids, err := repo.workExp.getIdsByTeacherId(ctx, client, d.ID)
			if err != nil {
				return err
			}
			requestIds := extractIdsFromWorkExp(d.Edges.WorkExperience)
			createable, updateable, deleteable := makeDataForCondition(requestIds, ids)

			if len(createable) > 0 {
				createData := filterWorkExperience(d.Edges.WorkExperience, createable)
				err = repo.workExp.createMany(ctx, client, createData, d.ID)
				if err != nil {
					return err
				}
			}
			if len(updateable) > 0 {
				updateData := filterWorkExperience(d.Edges.WorkExperience, updateable)
				err = repo.workExp.updateMany(ctx, client, updateData, d.ID)
				if err != nil {
					return err
				}
			}

			if len(deleteable) > 0 {
				_, err = repo.workExp.deletesByIds(ctx, client, deleteable, d.ID)
				if err != nil {
					return err
				}
			}
		} else {
			repo.workExp.deletesByTecaherId(ctx, client, d.ID)
		}

		// Certification
		if len(d.Edges.Certification) > 0 {
			ids, err := repo.certifi.getIdsByTeacherId(ctx, client, d.ID)
			if err != nil {
				return err
			}
			requestIds := extractIdsFromCertifi(d.Edges.Certification)
			createable, updateable, deleteable := makeDataForCondition(requestIds, ids)

			if len(createable) > 0 {
				createData := filterCertification(d.Edges.Certification, createable)
				err = repo.certifi.createMany(ctx, client, createData, d.ID)
				if err != nil {
					return err
				}
			}
			if len(updateable) > 0 {
				updateData := filterCertification(d.Edges.Certification, updateable)
				err = repo.certifi.updateMany(ctx, client, updateData, d.ID)
				if err != nil {
					return err
				}
			}

			if len(deleteable) > 0 {
				_, err = repo.certifi.deletesByIds(ctx, client, deleteable, d.ID)
				if err != nil {
					return err
				}
			}
		} else {
			// do
			repo.certifi.deletebyTeacherId(ctx, client, d.ID)
		}

		// YogaRaw
		if len(d.Edges.YogaRaw) > 0 {
			ids, err := repo.yoga.getRawIdsByTeacherId(ctx, client, d.ID)
			if err != nil {
				return err
			}

			requestIds := extractIdsFromYogaRaws(d.Edges.YogaRaw)
			createable, updateable, deleteable := makeDataForCondition(requestIds, ids)

			var teacherId *int
			if d.ID != 0 {
				teacherId = &d.ID
			}

			if len(createable) > 0 {
				createData := filterYogaRaws(d.Edges.YogaRaw, createable)
				err = repo.yoga.createRaws(ctx, client, createData, teacherId, nil)
				if err != nil {
					return err
				}
			}
			if len(updateable) > 0 {
				updateData := filterYogaRaws(d.Edges.YogaRaw, updateable)
				err = repo.yoga.updateRaws(ctx, client, updateData, teacherId, nil)
				if err != nil {
					return err
				}
			}

			if len(deleteable) > 0 {
				_, err = repo.yoga.deleteRawsByIds(ctx, client, deleteable, teacherId, nil)
				if err != nil {
					return err
				}
			}

		} else {
			// do
			repo.yoga.deleteRawsbyTeacherId(ctx, client, d.ID)
		}

		return
	})
}

func (repo *teacherRepository) Patch(ctx context.Context, d *request.TeacherPatchBody, id, userId int) (err error) {
	teacherInfo := structs.New(d.TeacherInfo)
	updateableTeacherInfo := utils.GetUpdateableDataV2(teacherInfo, teacher.Columns)

	return entx.WithTx(ctx, repo.db.Debug(), func(tx *ent.Tx) (err error) {
		client := tx.Client()

		clauseTeacher := client.Teacher.Update().
			Where(teacher.IDEQ(id), teacher.UserIDEQ(userId))

		for key, val := range updateableTeacherInfo {
			clauseTeacher.Mutation().SetField(key, val)
		}

		if err := clauseTeacher.Exec(ctx); err != nil {
			return err
		}

		if d.WorkExperiences != nil {
			for _, v := range d.WorkExperiences {
				s := structs.New(v)
				c := client.TeacherWorkExperience
				res := utils.GetUpdateableDataV2(s, twe.Columns)

				// Update
				if v.Id != nil {
					u := c.Update().Where(twe.TeacherIDEQ(id), twe.IDEQ(*v.Id))

					for key, val := range res {
						u.Mutation().SetField(key, val)
					}

					if err = u.Exec(ctx); err != nil {
						return
					}

					// Create
				} else {
					cr := c.Create().SetTeacherID(id)

					for key, val := range res {
						cr.Mutation().SetField(key, val)
					}
					if err = cr.Exec(ctx); err != nil {
						return
					}

				}

			}
		}

		return err
	})
}

func (repo *teacherRepository) GetOnlyIdByUserId(ctx context.Context, userId int) (id int, err error) {
	return repo.db.Teacher.Query().Where(teacher.UserIDEQ(userId)).OnlyID(ctx)
}

func (repo *teacherRepository) Exist(ctx context.Context, id int) (bool, error) {
	return repo.db.Teacher.Query().Where(teacher.IDEQ(id)).Exist(ctx)
}

func makeDataForCondition(requestIds []int, existIds []int) (createable []int, updateable []int, deleteable []int) {
	updateable = utils.Intersection(requestIds, existIds)
	deleteable = utils.Difference(existIds, requestIds)
	createable = utils.Difference(requestIds, existIds)
	return
}

func extractIdsFromWorkExp(val []*ent.TeacherWorkExperience) []int {
	var result []int
	for _, s := range val {
		result = append(result, s.ID)
	}
	return result
}

func filterWorkExperience(val []*ent.TeacherWorkExperience, ids []int) []*ent.TeacherWorkExperience {
	result := make([]*ent.TeacherWorkExperience, 0)

	for k := 0; k < len(ids); k++ {
		for i := k; i < len(val); i++ {
			if val[i].ID == ids[k] {
				result = append(result, val[i])
				break
			}
		}
	}

	return result
}

func extractIdsFromCertifi(val []*ent.TeacherCertification) []int {
	var result []int
	for _, s := range val {
		result = append(result, s.ID)
	}
	return result
}

func filterCertification(val []*ent.TeacherCertification, ids []int) []*ent.TeacherCertification {
	result := make([]*ent.TeacherCertification, 0)

	for k := 0; k < len(ids); k++ {
		for i := k; i < len(val); i++ {
			if val[i].ID == ids[k] {
				result = append(result, val[i])
				break
			}
		}
	}
	return result
}

func extractIdsFromYogaRaws(val []*ent.YogaRaw) []int {
	var result []int
	for _, s := range val {
		result = append(result, s.ID)
	}
	return result
}

func filterYogaRaws(val []*ent.YogaRaw, ids []int) []*ent.YogaRaw {
	result := make([]*ent.YogaRaw, 0)
	for k := 0; k < len(ids); k++ {
		for i := k; i < len(val); i++ {
			if val[i].ID == ids[k] {
				result = append(result, val[i])
				break
			}
		}
	}
	return result
}
