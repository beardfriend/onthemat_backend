package repository

import (
	"context"

	"onthemat/internal/app/model"
	"onthemat/internal/app/transport"
	"onthemat/internal/app/transport/request"
	"onthemat/internal/app/utils"
	"onthemat/pkg/ent"
	"onthemat/pkg/ent/teacher"
	"onthemat/pkg/ent/teacherworkexperience"
	"onthemat/pkg/ent/user"

	"onthemat/pkg/entx"

	"github.com/fatih/structs"
	"github.com/iancoleman/strcase"
)

type TeacherRepository interface {
	Create(ctx context.Context, d *ent.Teacher) error
	Update(ctx context.Context, d *ent.Teacher) (err error)
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

		if d.Edges.WorkExperience == nil {
			repo.workExp.deletesByTecaherId(ctx, client, d.ID)
		} else if len(d.Edges.WorkExperience) > 0 {
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

func (repo *teacherRepository) Patch(ctx context.Context, d *request.TeacherPatchBody, id, userId int) (err error) {
	updateableTeacherInfo := utils.GetUpdateableData(d.TeacherInfo, teacher.Columns)

	return entx.WithTx(ctx, repo.db.Debug(), func(tx *ent.Tx) (err error) {
		client := tx.Client()
		clauseTeacher := client.Teacher.Update().Where(teacher.IDEQ(id), teacher.UserIDEQ(userId))
		for key, val := range updateableTeacherInfo {
			clauseTeacher.Mutation().SetField(key, val)
		}

		if err := clauseTeacher.Exec(ctx); err != nil {
			return err
		}

		if d.WorkExperiences != nil {
			for _, v := range d.WorkExperiences {
				s := structs.New(v)
				if v.Id != nil {

					c := client.TeacherWorkExperience.Update().
						Where(
							teacherworkexperience.TeacherIDEQ(id),
							teacherworkexperience.IDEQ(*v.Id),
						)

					for _, f := range s.Fields() {
						switch f.Value().(type) {
						case *string:
							ptrValue := f.Value().(*string)
							if ptrValue != nil {
								value := *ptrValue
								c.Mutation().SetField(strcase.ToSnake(f.Name()), value)
							}
						case *int:
							ptrValue := f.Value().(*int)
							if ptrValue != nil {
								value := *ptrValue
								c.Mutation().SetField(strcase.ToSnake(f.Name()), value)
							}
						case *transport.TimeString:
							ptrValue := f.Value().(*transport.TimeString)
							if ptrValue != nil {
								value := *ptrValue
								c.Mutation().SetField(strcase.ToSnake(f.Name()), value)
							}
						}
					}

					err = c.Exec(ctx)
					if err != nil {
						return
					}

				} else {
					c := client.TeacherWorkExperience.Create().SetTeacherID(id)
					for _, f := range s.Fields() {
						switch f.Value().(type) {
						case *string:
							ptrValue := f.Value().(*string)
							if ptrValue != nil {
								value := *ptrValue
								c.Mutation().SetField(strcase.ToSnake(f.Name()), value)
							}
						case *int:
							ptrValue := f.Value().(*int)
							if ptrValue != nil {
								value := *ptrValue
								c.Mutation().SetField(strcase.ToSnake(f.Name()), value)
							}
						case *transport.TimeString:
							ptrValue := f.Value().(*transport.TimeString)
							if ptrValue != nil {
								value := *ptrValue
								c.Mutation().SetField(strcase.ToSnake(f.Name()), value)
							}
						}
					}
					err = c.Exec(ctx)
					if err != nil {
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
