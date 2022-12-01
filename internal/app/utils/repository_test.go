package utils_test

import (
	"fmt"
	"testing"
	"time"

	"onthemat/internal/app/transport"
	"onthemat/internal/app/utils"
	"onthemat/pkg/ent"
	"onthemat/pkg/ent/teacherworkexperience"

	"github.com/fatih/structs"
)

func TestGetUpdateableData(t *testing.T) {
	v := &ent.TeacherWorkExperience{
		ID:          1,
		TeacherID:   12,
		AcademyName: "academyName",
		Description: utils.String("ㅇㅇ"),
		WorkStartAt: transport.TimeString(time.Now()),
	}

	s := structs.New(v)
	asdf := utils.GetUpdateableDataV2(s, teacherworkexperience.Columns)
	fmt.Println(asdf)
}

func BenchmarkExtract(b *testing.B) {
	v := &ent.TeacherWorkExperience{
		ID:          1,
		TeacherID:   12,
		AcademyName: "academyName",
		Description: utils.String("ㅇㅇ"),
		WorkStartAt: transport.TimeString(time.Now()),
	}

	for i := 0; i < b.N; i++ {
		utils.GetUpdateableData(v, teacherworkexperience.Columns)
	}
}

func BenchmarkTwo(b *testing.B) {
	v := &ent.TeacherWorkExperience{
		ID:          1,
		TeacherID:   12,
		AcademyName: "academyName",
		Description: utils.String("ㅇㅇ"),
		WorkStartAt: transport.TimeString(time.Now()),
	}
	for i := 0; i < b.N; i++ {
		s := structs.New(v)
		utils.GetUpdateableDataV2(s, teacherworkexperience.Columns)
	}
}
