package model

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type TeacherWorkArea struct {
	ent.Schema
}

func (TeacherWorkArea) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "teacher_workarea"},
	}
}

func (TeacherWorkArea) Fields() []ent.Field {
	return []ent.Field{
		field.Int("teacher_id").
			Optional().
			Comment("foreignKey"),

		field.String("location").
			Optional().
			Nillable().
			Comment("지역"),

		field.String("gu").
			Comment("가능한 구역"),
	}
}

func (TeacherWorkArea) Mixin() []ent.Mixin {
	return []ent.Mixin{
		DefaultTimeMixin{},
	}
}

func (TeacherWorkArea) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("Teacher", Teacher.Type).
			Ref("workArea").
			Unique().Field("teacher_id"),
	}
}
