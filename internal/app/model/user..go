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
			Comment("휴대폰 번호"),

		field.Bool("is_super").
			Default(false).
			Comment("슈퍼 어드민인지 여부"),

		field.Bool("is_lock").
			Default(false).
			Comment("계정 잠겼는지 여부"),

		field.Int8("tries").Max(6).
			Default(0).
			Comment("로그인 시도 횟수"),

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
		edge.To("UserNormal", UserNormal.Type).
			StorageKey(edge.Column("id")),

		edge.To("Academy", Acadmey.Type).
			StorageKey(edge.Column("id")),

		edge.To("Teacher", Teacher.Type).
			StorageKey(edge.Column("id")),
	}
}
