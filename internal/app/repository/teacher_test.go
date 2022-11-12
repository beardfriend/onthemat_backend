package repository

import (
	"context"
	"fmt"
	"testing"

	"onthemat/internal/app/infrastructor"
	"onthemat/pkg/ent"
)

func TestTeacherCreate(t *testing.T) {
	db := infrastructor.NewPostgresDB()
	teacher := NewTeacherRepository(db)
	if err := teacher.Create(context.Background(), &ent.Teacher{
		Name: "name",
		Age:  26,
	}, 3); err != nil {
		t.Error(err)
	}
}

func TestTeacherList(t *testing.T) {
	db := infrastructor.NewPostgresDB()
	teacher := NewTeacherRepository(db)

	query := []string{"하타", "인"}
	var yogaSorts []*string
	for i := 0; i < len(query); i++ {
		v := &query[i]
		yogaSorts = append(yogaSorts, v)
	}

	querytwo := []string{"강서구"}
	var gu []*string
	for _, v := range querytwo {
		gu = append(gu, &v)
	}
	arr, _ := teacher.List(context.Background(), yogaSorts, gu)
	fmt.Println(arr)
}
