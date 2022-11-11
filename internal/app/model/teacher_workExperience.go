package model

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

type TeacherWorkExperience struct {
	ent.Schema
}

func (TeacherWorkExperience) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "teacher_workexperiences"},
	}
}

type ClassContent struct {
	YogaSort    string `json:"yogaSort"`
	RunningTime int    `json:"runningTime"`
}

func (TeacherWorkExperience) Fields() []ent.Field {
	return []ent.Field{
		field.String("academyName").
			Comment("근무지 이름"),

		field.String("image_url").
			Comment("자격증 사진"),

		field.Time("workStartAt").
			Comment("근무 시작일"),

		field.Time("workEndAt").
			Comment("근무 종료일"),

		field.Text("description").
			Comment("기타 설명"),

		field.JSON("classContent", &ClassContent{}).
			Comment("수업 내용"),
	}
}

func (TeacherWorkExperience) Mixin() []ent.Mixin {
	return []ent.Mixin{
		DefaultTimeMixin{},
	}
}
