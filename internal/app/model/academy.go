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
			NotEmpty().
			Comment("사업자 번호"),

		field.String("call_number").
			NotEmpty().
			Comment("학원 연락처"),

		field.String("address_road").
			NotEmpty().
			Comment("전체 주소"),

		field.String("address_sigun").
			NotEmpty().
			Comment("시 or 군"),

		field.String("address_gu").
			NotEmpty().
			Comment("구"),

		field.String("address_dong").
			NotEmpty().
			Comment("동"),

		field.String("address_detail").
			Optional().
			Comment("상세주소"),

		field.String("address_x").
			NotEmpty().
			Comment("x좌표"),

		field.String("address_y").
			NotEmpty().
			Comment("y좌표"),
	}
}

func (Acadmey) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("recruitment", Recruitment.Type),

		edge.From("user", User.Type).
			Ref("Academy").Unique().Required(),
	}
}
