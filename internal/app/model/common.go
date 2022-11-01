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
		field.Time("created_at").
			Immutable().
			Default(time.Now),
		field.Time("updated_at").
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
		field.Time("created_at").
			Immutable().
			Default(time.Now),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
		field.Time("removed_at").
			Optional().
			UpdateDefault(time.Now),
	}
}
