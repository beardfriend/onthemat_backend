package model

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// ------------------- Default  -------------------

type DefaultTimeMixin struct {
	mixin.Schema
}

func (DefaultTimeMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Time("createdAt").
			Immutable().
			Default(time.Now),
		field.Time("updatedAt").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// ------------------- WithRemovedAt  -------------------

type TimeWithRemovedAtMixin struct {
	mixin.Schema
}

func (TimeWithRemovedAtMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Time("createdAt").
			Immutable().
			Default(time.Now),
		field.Time("updatedAt").
			Default(time.Now).
			UpdateDefault(time.Now),
		field.Time("removedAt").
			Optional(),
	}
}
