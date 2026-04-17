package userDomain

type UserRole string

type UserStatus string

const (
	RoleAdmin UserRole = "ADMIN"
	RoleUser  UserRole = "USER"

	Active  UserStatus = "ACTIVE"
	Blocked UserStatus = "BLOCKED"
	//Deleted UserStatus = "DELETED"
)
