package mixin

import (
	entUtils "libs/ent-utils"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

type BaseMixin struct {
	ent.Schema
}

func (BaseMixin) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(entUtils.NewUUID).
			Unique().
			Immutable(),

		field.Time("created_at").
			Default(entUtils.NewTime).
			Immutable(),

		field.Time("updated_at").
			Default(entUtils.NewTime).
			UpdateDefault(entUtils.NewTime),

		field.Time("deleted_at").
			Optional().
			Nillable(),
	}
}

func (BaseMixin) Edges() []ent.Edge {
	return nil
}

func (BaseMixin) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("deleted_at"),
	}
}
