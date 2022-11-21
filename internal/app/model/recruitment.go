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
		field.Bool("isFinish").
			Default(false).
			Comment("채용 종료 여부"),
	}
}

func (Recruitment) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("instead", RecruitmentInstead.Type).
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),

		edge.From("writer", Academy.Type).
			Ref("recruitment").
			Unique(),
	}
}
