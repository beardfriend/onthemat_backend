package model

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type YogaRaw struct {
	ent.Schema
}

func (YogaRaw) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "yoga_raw"},
	}
}

func (YogaRaw) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Comment("요가 이름"),

		field.Bool("is_migrated").
			Comment("정식 그룹에 추가됐는지"),
	}
}

func (YogaRaw) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("yoga_raw").Unique(),
	}
}

func (YogaRaw) Mixin() []ent.Mixin {
	return []ent.Mixin{
		DefaultTimeMixin{},
	}
}
