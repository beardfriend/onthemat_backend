package model

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Acadmey struct {
	ent.Schema
}

func (Acadmey) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "academis"},
	}
}

func (Acadmey) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").SchemaType(map[string]string{
			dialect.Postgres: "varchar(30)",
		}).
			MaxLen(30).
			NotEmpty().
			Comment("학원 이름"),
	}
}

func (Acadmey) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("Academy").Unique().Required(),
	}
}
