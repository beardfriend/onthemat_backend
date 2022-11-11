package model

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type RecruitmentInstead struct {
	ent.Schema
}

func (RecruitmentInstead) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "recruitment_instead"},
	}
}

func (RecruitmentInstead) Fields() []ent.Field {
	return []ent.Field{
		field.String("yogaSort").
			Comment("요가 종류"),

		field.String("minCareer").
			Comment("최소 경력"),

		field.String("pay").
			Comment("급여"),

		field.Time("startDateTime").
			Comment("수업 시작 일시"),

		field.Time("endDateTime").
			Comment("수업 종료 일시"),
	}
}

func (RecruitmentInstead) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("applicant", Teacher.Type),
	}
}
