package userMapper

import (
	"github.com/deniSSTK/task-engine/auth-service/ent"
	userDomain "github.com/deniSSTK/task-engine/libs/user"
)

func MapEntUserToDomain(rawUser *ent.User) *userDomain.User {
	user := userDomain.User{
		Id:         rawUser.ID,
		Name:       rawUser.Name,
		SecondName: &rawUser.SecondName,
		Email:      rawUser.Email,
		Role:       userDomain.UserRole(rawUser.Role),
		Status:     userDomain.UserStatus(rawUser.Status),
	}

	return userDomain.NewUser(&user)
}
