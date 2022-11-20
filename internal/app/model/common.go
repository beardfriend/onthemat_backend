package model

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
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
			Immutable().SchemaType(map[string]string{
			dialect.Postgres: "timestamp",
		}).
			Default(time.Now),

		field.Time("updatedAt").
			Default(time.Now).
			Immutable().SchemaType(map[string]string{
			dialect.Postgres: "timestamp",
		}).
			UpdateDefault(time.Now),
	}
}

// ------------------- WithRemovedAt  -------------------
