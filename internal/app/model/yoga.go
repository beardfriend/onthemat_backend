package model

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Yoga struct {
	ent.Schema
}

func (Yoga) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "yoga"},
	}
}

func (Yoga) Fields() []ent.Field {
	return []ent.Field{
		field.String("name_kor").
			Comment("요가 이름 한국어"),

		field.String("name_eng").
			Optional().
			Nillable().
			Comment("요가 이름 영어"),

		field.Int("level").Optional().Nillable(),

		field.Text("description").
			Optional().
			Nillable().
			Comment("요가에 대한 설명"),
	}
}

func (Yoga) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("yoga"),

		edge.From("yoga_group", YogaGroup.Type).
			Ref("yoga").
			Unique(),
	}
}

func (Yoga) Mixin() []ent.Mixin {
	return []ent.Mixin{
		DefaultTimeMixin{},
	}
}
