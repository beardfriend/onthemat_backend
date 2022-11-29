package model

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type AreaSiGungu struct {
	ent.Schema
}

func (AreaSiGungu) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "areas_sigungu"},
	}
}

func (AreaSiGungu) Fields() []ent.Field {
	return []ent.Field{
		field.Int("area_sido_id").
			Optional(),

		field.String("name").
			Comment("시군구 이름"),

		field.String("adm_code").
			Comment("시군구 코드").
			Unique(),

		field.String("parent_code").
			Optional().
			Nillable().
			Comment("수원시 장안구 같은 경우에는 수원시에 해당"),

		field.Int32("version").
			Comment("정보 버전"),
	}
}

func (AreaSiGungu) Edges() []ent.Edge {
	return []ent.Edge{
		// Many To One
		edge.From("sido", AreaSiDo.Type).
			Ref("sigungu").
			Unique().
			Field("area_sido_id"),

		edge.To("academy", Academy.Type),

		edge.From("teacher", Teacher.Type).
			Ref("sigungu"),
	}
}
