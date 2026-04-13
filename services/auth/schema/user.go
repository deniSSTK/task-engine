package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/deniSSTK/task-engine/libs/mixin"
	userDomain "github.com/deniSSTK/task-engine/libs/user"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.BaseMixin{},
	}
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			MinLen(1).
			MaxLen(100),

		field.String("second_name").
			MinLen(1).
			MaxLen(100).
			Optional().
			NotEmpty(),

		field.String("full_name").
			MaxLen(201),

		field.String("email").
			MinLen(5).
			MaxLen(200).
			Unique(),

		field.String("password_hash"),

		field.Time("last_login_at").
			Optional().
			Nillable(),

		field.Enum("role").
			Values(
				string(userDomain.Admin),
				string(userDomain.User),
			).
			Default(string(userDomain.User)),

		field.Enum("status").
			Values(
				string(userDomain.Active),
				string(userDomain.Blocked),
			).
			Default(string(userDomain.Active)),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}
