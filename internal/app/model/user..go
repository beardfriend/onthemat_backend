package model

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type User struct {
	ent.Schema
}

func (User) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "users"},
	}
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("email").SchemaType(map[string]string{
			dialect.Postgres: "varchar(100)",
		}).MaxLen(100).
			NotEmpty().
			Comment("이메일"),

		field.String("password").
			Optional().
			Sensitive().
			Comment("패스워드"),

		field.String("social_name").
			Optional().
			Comment("소셜 로그인을 제공한 업체 이름"),

		field.String("social_key").Optional().
			Comment("소셜 로그인 시 발급되는 고유 키"),

		field.String("nickname").SchemaType(map[string]string{
			dialect.Postgres: "varchar(20)",
		}).
			MaxLen(20).
			NotEmpty().
			Comment("닉네임"),

		field.String("phone_num").
			SchemaType(map[string]string{
				dialect.Postgres: "varchar(11)",
			}).
			MaxLen(11).
			Optional().
			Comment("휴대폰 번호"),

		field.Time("term_agree_at").
			Default(time.Now()).
			Comment("약관 동의 일시"),

		field.Time("last_login_at").
			Default(time.Now()).
			Comment("마지막 로그인 일시"),
	}
}

func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeWithRemovedAtMixin{},
	}
}

func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("Academy", Acadmey.Type).
			StorageKey(edge.Column("id")),

		edge.To("Teacher", Teacher.Type).
			StorageKey(edge.Column("id")),
	}
}
