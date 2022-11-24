package model

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type AreaSiDo struct {
	ent.Schema
}

func (AreaSiDo) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "areas_sido"},
	}
}

func (AreaSiDo) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Comment("시도 이름"),

		field.String("adm_code").
			Comment("행정구역 코드").
			Unique(),

		field.Int32("version").
			Comment("정보 버전"),
	}
}

func (AreaSiDo) Edges() []ent.Edge {
	return []ent.Edge{
		// One to Many
		edge.To("sigungu", AreaSiGungu.Type),
	}
}
