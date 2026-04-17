package userDomain

import "github.com/google/uuid"

type User struct {
	Id         uuid.UUID
	Name       string
	SecondName *string
	Email      string
	Role       UserRole
	Status     UserStatus
	FullName   string
}

func (u *User) BuildFullName() string {
	if u.SecondName == nil || *u.SecondName == "" {
		return u.Name
	}

	return u.Name + " " + *u.SecondName
}

func NewUser(
	user *User,
) *User {
	user.FullName = user.BuildFullName()

	return user
}
