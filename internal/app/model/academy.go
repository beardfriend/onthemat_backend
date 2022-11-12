package model

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Academy struct {
	ent.Schema
}

func (Academy) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "academies"},
	}
}

func (Academy) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").SchemaType(map[string]string{
			dialect.Postgres: "varchar(30)",
		}).
			MaxLen(30).
			NotEmpty().
			Comment("학원 이름"),

		field.String("businessCode").SchemaType(map[string]string{
			dialect.Postgres: "varchar(30)",
		}).
			MaxLen(15).
			NotEmpty().
			Comment("사업자 번호"),

		field.String("callNumber").
			NotEmpty().
			Comment("학원 연락처"),

		field.String("addressRoad").
			NotEmpty().
			Comment("전체 주소"),

		field.String("addressSigun").
			NotEmpty().
			Comment("시 or 군"),

		field.String("addressGu").
			NotEmpty().
			Comment("구"),

		field.String("addressDong").
			NotEmpty().
			Comment("동"),

		field.String("addressDetail").
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

func (Academy) Mixin() []ent.Mixin {
	return []ent.Mixin{
		DefaultTimeMixin{},
	}
}

func (Academy) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("recruitment", Recruitment.Type),

		edge.From("user", User.Type).
			Ref("Academy").Unique().Required(),
	}
}
