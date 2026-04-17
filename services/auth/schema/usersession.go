package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/deniSSTK/task-engine/libs/mixin"
	"github.com/google/uuid"
)

// UserSession holds the schema definition for the UserSession entity.
type UserSession struct {
	ent.Schema
}

func (UserSession) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.BaseMixin{},
	}
}

// Fields of the UserSession.
func (UserSession) Fields() []ent.Field {
	return []ent.Field{
		field.String("refresh_token").
			NotEmpty().
			Immutable(),

		field.String("ip").
			Optional().
			Nillable().
			Immutable().
			NotEmpty(),

		field.String("user_agent").
			Optional().
			Nillable().
			Immutable().
			NotEmpty(),

		field.Time("expires_at").
			Immutable(),

		field.UUID("user_id", uuid.UUID{}).
			Immutable(),

		//TODO: add device_id field
		//field.UUID("device_id", uuid.UUID{}).
		//	Default(uuid.New),
	}
}

// Edges of the UserSession.
func (UserSession) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Field("user_id").
			Ref("sessions").
			Required().
			Unique().
			Immutable(),
	}
}
