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
		field.String("name").MaxLen(100),
		field.String("second_name").MaxLen(100).Optional(),
		field.String("full_name").MaxLen(201),

		field.String("email").Unique(),
		field.String("password"),

		field.Time("last_logined_at"),

		field.Enum("role").
			Values(
				string(user.Admin),
				string(user.User),
			).
			Default(string(user.User)),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}
