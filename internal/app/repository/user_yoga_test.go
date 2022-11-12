package repository

import (
	"context"
	"testing"

	"onthemat/internal/app/infrastructor"
	"onthemat/pkg/ent"
	"onthemat/pkg/ent/useryoga"
)

func TestCreateMany(t *testing.T) {
	context := context.Background()
	db := infrastructor.NewPostgresDB()
	repo := NewUserYogaRepository(db)

	var yoga []*ent.UserYoga
	yoga = append(yoga, &ent.UserYoga{
		Name:     "하타",
		UserType: useryoga.UserTypeTeacher,
	})

	yoga = append(yoga, &ent.UserYoga{
		Name:     "인",
		UserType: useryoga.UserTypeTeacher,
	})
	if err := repo.CreateMany(context, yoga, 3); err != nil {
		if ent.IsConstraintError(err) {
			t.Error("중복된 값이 존재합니다.")
			return
		}
		t.Error(err)
		return
	}
}
