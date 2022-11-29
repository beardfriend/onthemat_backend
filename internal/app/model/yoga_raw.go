package model

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
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
		field.Int("academy_id").
			Optional().
			Nillable(),

		field.Int("teacher_id").
			Optional().
			Nillable(),

		field.String("name").
			Comment("요가 이름"),

		field.Bool("is_migrated").
			Default(false).
			Comment("정식 그룹에 추가됐는지"),
	}
}

func (YogaRaw) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name", "academy_id").
			Unique(),

		index.Fields("name", "teacher_id").
			Unique(),
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
