package schema

import (
	"libs/mixin"
	"libs/user"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
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
			MaxLen(100),

		field.String("second_name").MaxLen(100).
			Optional().
			NotEmpty(),

		field.String("full_name").
			MaxLen(201),

		field.String("email").
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
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}
