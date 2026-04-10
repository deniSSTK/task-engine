package userDomain

type UserRole string

type UserStatus string

const (
	Admin UserRole = "ADMIN"
	User  UserRole = "USER"

	Active  UserStatus = "ACTIVE"
	Blocked UserStatus = "BLOCKED"
)
