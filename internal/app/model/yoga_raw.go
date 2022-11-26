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
		field.Int("academy_id").Optional(),

		field.Int("teacher_id").Optional(),

		field.String("name").
			Comment("요가 이름"),

		field.Bool("is_migrated").
			Default(false).
			Comment("정식 그룹에 추가됐는지"),
	}
}

func (YogaRaw) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("academy", Academy.Type).
			Ref("yogaRaw").
			Unique().
			Field("academy_id"),

		edge.From("teacher", Teacher.Type).
			Ref("yogaRaw").
			Unique().
			Field("teacher_id"),
	}
}

func (YogaRaw) Mixin() []ent.Mixin {
	return []ent.Mixin{
		DefaultTimeMixin{},
	}
}
