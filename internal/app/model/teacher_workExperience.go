package model

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type TeacherWorkExperience struct {
	ent.Schema
}

func (TeacherWorkExperience) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "teacher_workexperiences"},
	}
}

func (TeacherWorkExperience) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id"),

		field.Int("teacher_id"),

		field.String("academyName").
			Comment("근무지 이름"),

		field.Time("workStartAt").
			SchemaType(
				map[string]string{
					dialect.Postgres: "timestamp",
				},
			).
			Comment("근무 시작일"),

		field.Time("workEndAt").
			SchemaType(
				map[string]string{
					dialect.Postgres: "timestamp",
				},
			).
			Optional().
			Nillable().
			Comment("근무 종료일"),

		field.Text("description").
			Optional().
			Nillable().
			Comment("기타 설명"),
	}
}

func (TeacherWorkExperience) Mixin() []ent.Mixin {
	return []ent.Mixin{
		DefaultTimeMixin{},
	}
}

func (TeacherWorkExperience) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("Teacher", Teacher.Type).
			Ref("workExperience").
			Unique().
			Required().
			Field("teacher_id"),
	}
}
