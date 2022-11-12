package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"onthemat/internal/app/infrastructor"
	"onthemat/internal/app/model"
	"onthemat/pkg/ent"
)

func TestCRUD(t *testing.T) {
	ctx := context.Background()
	db := infrastructor.NewPostgresDB()
	wexpModule := NewTeacherWorkExperience(db)
	des := "dd"

	firstObject := &ent.TeacherWorkExperience{
		AcademyName: "목동 하타요가",
		ImageURL:    "http://www.naver.com",
		WorkStartAt: time.Now(),
		WorkEndAt:   time.Now().Add(time.Hour * 3),
		Description: &des,
		ClassContent: &[]model.ClassContent{
			{
				YogaSort:    "하타",
				RunningTime: 20,
			},

			{
				YogaSort:    "아헹가",
				RunningTime: 200,
			},
		},
	}

	secondObject := &ent.TeacherWorkExperience{
		AcademyName: "목동 아리아요가",
		ImageURL:    "http://www.naver.com",
		WorkStartAt: time.Now(),
		WorkEndAt:   time.Now().Add(time.Hour * 3),
		Description: &des,
		ClassContent: &[]model.ClassContent{
			{
				YogaSort:    "아쉬탕가",
				RunningTime: 20,
			},

			{
				YogaSort:    "아헹가",
				RunningTime: 200,
			},
		},
	}
	wexpModule.Create(ctx, firstObject, 4)

	wexpModule.CreateMany(ctx, []*ent.TeacherWorkExperience{firstObject, secondObject}, 4)

	d, _ := wexpModule.Get(ctx, 1)
	fmt.Println(d)

	data, _ := wexpModule.ListByTeacherID(ctx, 4)
	fmt.Println(data)
}
