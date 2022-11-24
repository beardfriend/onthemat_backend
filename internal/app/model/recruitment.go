package model

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Recruitment struct {
	ent.Schema
}

func (Recruitment) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "recruitments"},
	}
}

func (Recruitment) Fields() []ent.Field {
	return []ent.Field{
		field.Int("academy_id"),

		field.Bool("isFinish").
			Default(false).
			Comment("채용 종료 여부"),
	}
}

func (Recruitment) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("recruitmentInstead", RecruitmentInstead.Type).
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),

		edge.To("yoga", Yoga.Type),

		edge.From("writer", Academy.Type).
			Ref("recruitment").
			Unique().
			Required().
			Field("academy_id"),
	}
}
