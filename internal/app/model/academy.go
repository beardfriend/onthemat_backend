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
		field.Int("id"),

		field.Int("user_id").
			StructTag(`json:"userId,omitempty"`).
			Comment("foreignKey"),

		field.Int("sigungu_id").
			StructTag(`json:"sigunguId,omitempty"`).
			Comment("foreignKey"),

		field.String("name").
			SchemaType(
				map[string]string{
					dialect.Postgres: "varchar(30)",
				},
			).
			NotEmpty().
			Comment("학원 이름"),

		field.String("businessCode").
			SchemaType(
				map[string]string{
					dialect.Postgres: "varchar(30)",
				},
			).
			NotEmpty().
			Comment("사업자 번호"),

		field.String("callNumber").
			NotEmpty().
			Comment("학원 연락처"),

		field.String("addressRoad").
			NotEmpty().
			Comment("전체 주소"),

		field.String("addressDetail").
			Optional().
			Nillable().
			Comment("상세주소"),
	}
}

func (Academy) Mixin() []ent.Mixin {
	return []ent.Mixin{
		DefaultTimeMixin{},
	}
}

func (Academy) Edges() []ent.Edge {
	return []ent.Edge{
		// 다루는 요가
		edge.To("yoga", Yoga.Type),

		edge.To("yogaRaw", YogaRaw.Type),

		edge.To("recruitment", Recruitment.Type).
			Annotations(
				entsql.Annotation{
					OnDelete: entsql.Cascade,
				}),

		edge.From("user", User.Type).
			Ref("Academy").
			Unique().
			Required().
			Field("user_id"),

		edge.From("area_sigungu", AreaSiGungu.Type).
			Ref("academy").
			Unique().
			Required().
			Field("sigungu_id"),
	}
}
