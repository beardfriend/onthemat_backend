package model

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

type TeacherCertification struct {
	ent.Schema
}

func (TeacherCertification) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "teacher_certifications"},
	}
}

func (TeacherCertification) Fields() []ent.Field {
	return []ent.Field{
		field.String("agencyName").
			Comment("자격증 기관명"),

		field.String("imageRrl").
			Comment("자격증 사진"),

		field.Time("classStartAt").
			Comment("자격증 수업 시작일"),

		field.Time("classEndAt").
			Comment("자격증 수업 종료일"),

		field.Text("description").
			Comment("기타 설명"),
	}
}

func (TeacherCertification) Mixin() []ent.Mixin {
	return []ent.Mixin{
		DefaultTimeMixin{},
	}
}
