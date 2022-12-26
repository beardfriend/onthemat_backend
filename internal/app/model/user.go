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
		field.String("email").
			SchemaType(map[string]string{
				dialect.Postgres: "varchar(100)",
			}).MaxLen(100).
			Optional().
			Nillable().
			Unique().
			Comment("이메일"),

		field.Bool("isEmailVerified").
			Default(false).
			Comment("이메일 인증 여부"),

		field.String("password").
			Optional().
			Nillable().
			Sensitive().
			Comment("패스워드"),

		field.String("tempPassword").
			Optional().
			Nillable().
			Sensitive().
			Comment("임시 비밀번호"),

		field.Int8("socialName").
			GoType(SocialType(0)).
			Optional().
			Nillable().
			Comment("소셜 로그인을 제공한 업체 이름 1:kakao, 2:google, 3:naver"),

		field.String("socialKey").
			Optional().
			Nillable().
			Unique().
			Comment("소셜 로그인 시 발급되는 고유 키"),

		field.String("nickname").
			SchemaType(map[string]string{
				dialect.Postgres: "varchar(50)",
			}).
			MaxLen(50).
			Optional().
			Nillable().
			Comment("닉네임"),

		field.Int8("type").
			GoType(UserType(0)).
			Optional().
			Nillable().
			Comment("유저 타입 1:teacher 2:academy"),

		field.String("phoneNum").
			SchemaType(map[string]string{
				dialect.Postgres: "varchar(11)",
			}).
			MaxLen(11).
			Optional().
			Nillable().
			Comment("휴대폰 번호"),

		field.Time("termAgreeAt").
			SchemaType(map[string]string{
				dialect.Postgres: "timestamp",
			}).
			Comment("약관 동의 일시"),

		field.Time("lastLoginAt").
			Default(time.Now).
			SchemaType(map[string]string{
				dialect.Postgres: "timestamp",
			}).
			Comment("마지막 로그인 일시"),
	}
}

func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		DefaultTimeMixin{},
	}
}

func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("Academy", Academy.Type).
			Unique().
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),

		edge.To("Teacher", Teacher.Type).
			Unique().
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),

		edge.To("Image", Image.Type),
	}
}
