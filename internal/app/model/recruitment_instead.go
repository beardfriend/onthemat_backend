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
		field.Int("recruitment_id"),

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
		edge.From("recuritment", Recruitment.Type).
			Ref("recruitmentInstead").
			Unique().
			Required().
			Field("recruitment_id"),

		edge.To("applicant", Teacher.Type).
			StorageKey(edge.Table("rinstead_academy"), edge.Columns("r_instead_id", "academy_id")),

		edge.To("yoga", Yoga.Type).
			StorageKey(edge.Table("rinstead_yoga"), edge.Columns("r_instead_id", "yoga_id")),
	}
}
