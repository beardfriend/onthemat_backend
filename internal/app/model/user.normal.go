package model

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type UserNormal struct {
	ent.Schema
}

func (UserNormal) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "users_normal"},
	}
}

func (UserNormal) Fields() []ent.Field {
	return []ent.Field{
		field.String("email").SchemaType(map[string]string{
			dialect.Postgres: "varchar(100)",
		}).
			MaxLen(100).
			NotEmpty().
			Comment("이메일"),

		field.String("password").
			NotEmpty().
			Sensitive().
			Comment("패스워드"),
	}
}

func (UserNormal) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("UserNormal").Unique().Required(),
	}
}
