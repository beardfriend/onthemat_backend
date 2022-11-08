package model

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Image struct {
	ent.Schema
}

func (Image) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "images"},
	}
}

func (Image) Fields() []ent.Field {
	return []ent.Field{
		field.String("path").
			Comment("이미지 주소"),

		field.Int("size").
			Comment("이미지 크기"),

		field.String("content_type").
			Comment("이미지 콘텐츠 타입"),

		field.String("name").
			Comment("이미지 이름"),

		field.Enum("type").
			Values("profile", "logo").
			Comment("이미지 타입"),
	}
}

func (Image) Mixin() []ent.Mixin {
	return []ent.Mixin{
		DefaultTimeMixin{},
	}
}

func (Image) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("Image").
			Unique(),
	}
}
