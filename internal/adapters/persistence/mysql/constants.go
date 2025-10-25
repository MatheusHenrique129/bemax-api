package mysql

import "errors"

var (
	ErrQuery                = errors.New("error while executing query")
	ErrRoleNotFound         = errors.New("role not found")
	ErrUserNotFound         = errors.New("user not found")
	ErrLoginNotFound        = errors.New("login not found")
	ErrNoRowsAffected       = errors.New("no rows affected")
	ErrRolesForUserNotFound = errors.New("no roles found for user")
)
