package model

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

type UserYoga struct {
	ent.Schema
}

func (UserYoga) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "user_yoga"},
	}
}

func (UserYoga) Fields() []ent.Field {
	return []ent.Field{
		field.Int("user_id"),

		field.String("name").
			Comment("이름"),

		field.Enum("userType").
			Values("teacher", "academy").
			Comment("유저 종류"),
	}
}

func (UserYoga) Edges() []ent.Edge {
	return []ent.Edge{}
}
