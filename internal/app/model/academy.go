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

		field.String("business_code").SchemaType(map[string]string{
			dialect.Postgres: "varchar(30)",
		}).
			MaxLen(15).
			Nillable().
			Comment("사업자 번호"),

		field.String("full_address").
			Optional().
			Comment("주소"),

		field.String("si").
			Optional().
			Comment("시"),

		field.String("gun").
			Optional().
			Comment("군"),

		field.String("gu").
			Optional().
			Comment("구"),

		field.String("dong").
			Optional().
			Comment("동"),

		field.String("x").
			Optional().
			Comment("x좌표"),

		field.String("y").
			Optional().
			Comment("y좌표"),
	}
}

func (Acadmey) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("Academy").Unique().Required(),
	}
}
