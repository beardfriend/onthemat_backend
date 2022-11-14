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
			Comment("큰 범주의 요가 이름 (한글)"),

		field.String("name_kor").
			Comment("요가 이름 한국어"),

		field.String("name_eng").
			Comment("요가 이름 영어"),

		field.Text("description").
			Optional().
			Comment("요가에 대한 설명"),
	}
}

func (Yoga) Mixin() []ent.Mixin {
	return []ent.Mixin{
		DefaultTimeMixin{},
	}
}
