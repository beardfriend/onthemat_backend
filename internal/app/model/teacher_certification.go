package model

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type TeacherCertification struct {
	ent.Schema
}

func (TeacherCertification) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "teacher_certifications"},
	}
}

func (TeacherCertification) Fields() []ent.Field {
	return []ent.Field{
		field.Int("teacher_id").
			Comment("foreignKey"),

		field.String("agencyName").
			Comment("자격증 기관명"),

		field.String("imageUrl").
			Optional().
			Nillable().
			Comment("자격증 사진"),

		field.Time("classStartAt").
			SchemaType(
				map[string]string{
					dialect.Postgres: "timestamp",
				},
			).
			Comment("자격증 수업 시작일"),

		field.Time("classEndAt").
			SchemaType(
				map[string]string{
					dialect.Postgres: "timestamp",
				},
			).
			Optional().
			Nillable().
			Comment("자격증 수업 종료일"),

		field.Text("description").
			Comment("기타 설명"),
	}
}

func (TeacherCertification) Mixin() []ent.Mixin {
	return []ent.Mixin{
		DefaultTimeMixin{},
	}
}

func (TeacherCertification) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("Teacher", Teacher.Type).
			Ref("certification").
			Unique().
			Required().
			Field("teacher_id"),
	}
}
