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
		field.Int("user_id").Comment("foreignKey"),

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
			Ref("Teacher").
			Unique().
			Required().
			Field("user_id"),

		edge.From("recruitment_instead", RecruitmentInstead.Type).
			Ref("applicant"),

		// 다루는 요가
		edge.To("yoga", Yoga.Type),

		edge.To("yogaRaw", YogaRaw.Type),

		// 자격증
		edge.To("certification", TeacherCertification.Type).
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),

		// 근무 경험
		edge.To("workExperience", TeacherWorkExperience.Type).
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),

		// 근무지
		edge.To("workArea", TeacherWorkArea.Type).
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),
	}
}
