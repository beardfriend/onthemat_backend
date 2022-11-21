package model

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type YogaGroup struct {
	ent.Schema
}

func (YogaGroup) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "yoga_group"},
	}
}

func (YogaGroup) Fields() []ent.Field {
	return []ent.Field{
		field.String("category").
			Comment("큰 범주의 요가 이름 (한글)"),

		field.String("category_eng").
			Comment("요가 이름 영어"),

		field.Text("description").
			Optional().
			Comment("요가에 대한 설명"),
	}
}

func (YogaGroup) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("yoga", Yoga.Type),
	}
}

func (YogaGroup) Mixin() []ent.Mixin {
	return []ent.Mixin{
		DefaultTimeMixin{},
	}
}
