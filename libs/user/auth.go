package userDomain

import "github.com/google/uuid"

type AuthUser struct {
	Id   uuid.UUID
	Role UserRole
}
