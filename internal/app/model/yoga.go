package model

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

type Yoga struct {
	ent.Schema
}

func (Yoga) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "yoga"},
	}
}

func (Yoga) Fields() []ent.Field {
	return []ent.Field{
		field.String("category").
			Comment("메인으로 묶이는 카테고리 ex) 아쉬탕가"),

		field.String("name").
			Comment("요가 이름"),

		field.Text("description").
			Optional().
			Comment("요가에 대한 설명"),

		field.Bool("is_offical").
			Comment("관리자가 등록했는지 여부"),
	}
}

func (Yoga) Mixin() []ent.Mixin {
	return []ent.Mixin{
		DefaultTimeMixin{},
	}
}
