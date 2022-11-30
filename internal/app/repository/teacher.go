package repository

import (
	"context"

	"onthemat/internal/app/model"
	"onthemat/pkg/ent"
	"onthemat/pkg/ent/teacher"
	"onthemat/pkg/ent/user"

	"onthemat/pkg/entx"
)

type TeacherRepository interface {
	Create(ctx context.Context, d *ent.Teacher) error
	Update(ctx context.Context, d *ent.Teacher) (err error)
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

		if len(d.Edges.WorkExperience) > 0 {
			ids, err := repo.workExp.getIdsByTeacherId(ctx, client, d.ID)
			if err != nil {
				return err
			}
			updateable := make([]*ent.TeacherWorkExperience, 0)
			createable := make([]*ent.TeacherWorkExperience, 0)
			createable = append(createable, d.Edges.WorkExperience...)

			for _, w := range d.Edges.WorkExperience {
				for _, id := range ids {
					if w.ID == id {
						updateable = append(updateable, w)
						for k, v := range createable {
							if v.ID == id {
								createable = append(createable[:k], createable[k+1:]...)
								break
							}
						}
						for k := 0; k < len(ids); k++ {
							if ids[k] == id {
								ids = append(ids[:k], ids[k+1:]...)
								break
							}
						}

						break
					}
				}
			}
			deleteable := ids

			if len(createable) > 0 {
				err = repo.workExp.createMany(ctx, client, createable, d.ID)
				if err != nil {
					return err
				}
			}
			if len(updateable) > 0 {
				err = repo.workExp.updateMany(ctx, client, updateable, d.ID)
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
		}

		if len(d.Edges.Certification) > 0 {
			ids, err := repo.certifi.getIdsByTeacherId(ctx, client, d.ID)
			if err != nil {
				return err
			}
			updateable := make([]*ent.TeacherCertification, 0)
			createable := make([]*ent.TeacherCertification, 0)
			createable = append(createable, d.Edges.Certification...)

			for _, w := range d.Edges.Certification {
				for _, id := range ids {
					if w.ID == id {
						updateable = append(updateable, w)

						for k, v := range createable {
							if v.ID == id {
								createable = append(createable[:k], createable[k+1:]...)
								break
							}
						}
						for k := 0; k < len(ids); k++ {
							if ids[k] == id {
								ids = append(ids[:k], ids[k+1:]...)
								break
							}
						}
						break
					}
				}
			}

			deleteable := ids

			if len(createable) > 0 {
				err = repo.certifi.createMany(ctx, client, createable, d.ID)
				if err != nil {
					return err
				}
			}
			if len(updateable) > 0 {
				err = repo.certifi.updateMany(ctx, client, updateable, d.ID)
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
		}

		if len(d.Edges.YogaRaw) > 0 {
			ids, err := repo.yoga.getRawIdsByTeacherId(ctx, client, d.ID)
			if err != nil {
				return err
			}
			updateable := make([]*ent.YogaRaw, 0)
			createable := make([]*ent.YogaRaw, 0)
			createable = append(createable, d.Edges.YogaRaw...)

			for _, w := range d.Edges.YogaRaw {
				for _, id := range ids {
					if w.ID == id {
						updateable = append(updateable, w)

						for k, v := range createable {
							if v.ID == id {
								createable = append(createable[:k], createable[k+1:]...)
								break
							}
						}
						for k := 0; k < len(ids); k++ {
							if ids[k] == id {
								ids = append(ids[:k], ids[k+1:]...)
								break
							}
						}
						break
					}
				}
			}

			deleteable := ids
			var teacherId *int
			if d.ID != 0 {
				teacherId = &d.ID
			}
			if len(createable) > 0 {

				err = repo.yoga.createRaws(ctx, client, createable, teacherId, nil)
				if err != nil {
					return err
				}
			}
			if len(updateable) > 0 {
				err = repo.yoga.updateRaws(ctx, client, updateable, teacherId, nil)
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

		}

		return
	})
}

func (repo *teacherRepository) GetOnlyIdByUserId(ctx context.Context, userId int) (id int, err error) {
	return repo.db.Teacher.Query().Where(teacher.UserIDEQ(userId)).OnlyID(ctx)
}
