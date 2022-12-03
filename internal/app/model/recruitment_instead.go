package model

import (
	"onthemat/internal/app/transport"

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

type Schedule struct {
	StartDateTime transport.TimeString `json:"startDateTime"`
	EndDateTime   transport.TimeString `json:"endDateTime"`
}

func (RecruitmentInstead) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id"),

		field.Int("recruitment_id"),

		field.Int("teacher_id").
			Optional().
			Nillable(),

		field.String("minCareer").
			Comment("최소 경력"),

		field.String("pay").
			Comment("급여"),

		field.JSON("schedule", []*Schedule{}).Optional(),
	}
}

func (RecruitmentInstead) Mixin() []ent.Mixin {
	return []ent.Mixin{
		DefaultTimeMixin{},
	}
}

func (RecruitmentInstead) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("recuritment", Recruitment.Type).
			Ref("recruitmentInstead").
			Unique().
			Required().
			Field("recruitment_id"),

		edge.From("passer", Teacher.Type).
			Ref("passer").
			Unique().
			Field("teacher_id"),

		edge.To("applicant", Teacher.Type).
			StorageKey(edge.Table("rinstead_teacher"), edge.Columns("r_instead_id", "teacher_id")),

		edge.To("yoga", Yoga.Type).
			StorageKey(edge.Table("rinstead_yoga"), edge.Columns("r_instead_id", "yoga_id")),
	}
}
