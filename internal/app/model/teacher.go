package model

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Teacher struct {
	ent.Schema
}

func (Teacher) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "teachers"},
	}
}

func (Teacher) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").SchemaType(map[string]string{
			dialect.Postgres: "varchar(10)",
		}).
			MaxLen(10).
			NotEmpty().
			Comment("선생님 이름"),

		field.Int("age").
			Comment("나이"),

		field.Bool("isProfileOpen").
			Default(false).
			Comment("프로필 오픈 여부"),
	}
}

func (Teacher) Mixin() []ent.Mixin {
	return []ent.Mixin{
		DefaultTimeMixin{},
	}
}

func (Teacher) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("Teacher").Unique().Required(),

		edge.From("recruitment_instead", RecruitmentInstead.Type).
			Ref("applicant"),

		// 자격증
		edge.To("certification", TeacherCertification.Type),

		// 근무 경험
		edge.To("workExperience", TeacherWorkExperience.Type),

		// 근무지
		edge.To("workArea", TeacherWorkArea.Type),
	}
}
