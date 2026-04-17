package userDomain

import "errors"

var (
	UserBlocked = errors.New("user blocked")
	UserDeleted = errors.New("user deleted")
)
