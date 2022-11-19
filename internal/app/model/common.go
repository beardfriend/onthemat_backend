package model

import (
	"database/sql/driver"
	"fmt"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// ------------------- Default  -------------------

type DefaultTimeMixin struct {
	mixin.Schema
}

type TimeString time.Time

func (t TimeString) Value() (driver.Value, error) {
	h := time.Time(t).Hour()

	min := time.Time(t).Minute()
	s := time.Time(t).Second()
	y, m, d := time.Time(t).Date()
	return time.Date(y, m, d, h, min, s, 0, time.Time(t).Location()), nil
}

func (t *TimeString) Scan(value any) error {
	switch value.(type) {
	case time.Time:
		*t = TimeString(value.(time.Time))
		break
	default:
		parsed, _ := time.Parse("2006-01-02 15:04:05", fmt.Sprintf("%s", value))
		*t = TimeString(parsed)
		break
	}
	return nil
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
