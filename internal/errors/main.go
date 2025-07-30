package errors

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type Error struct {
	Id   string
	Code int
	Err  error
}

func (e *Error) Error() string {
	return fmt.Sprintf("status %d: err %v", e.Code, e.Err)
}

func INTERNAL_SERVER_ERROR() *Error {
	return &Error{
		Id:   uuid.NewString(),
		Code: 500,
		Err:  errors.New("internal server error"),
	}
}

func PASSWORD_HASHING() *Error {
	return &Error{
		Id:   uuid.NewString(),
		Code: 500,
		Err:  errors.New("not possible to hash password"),
	}
}

func USERNAME_INVALID() *Error {
	return &Error{
		Id:   uuid.NewString(),
		Code: 400,
		Err:  errors.New("username is invalid. please check the requirements"),
	}
}

func EMAIL_INVALID() *Error {
	return &Error{
		Id:   uuid.NewString(),
		Code: 400,
		Err:  errors.New("email is invalid. please check the requirements"),
	}
}

func PASSWORD_INVALID() *Error {
	return &Error{
		Id:   uuid.NewString(),
		Code: 400,
		Err:  errors.New("password is invalid. please check the requirements"),
	}
}

func USER_NOT_FOUND_ERROR() *Error {
	return &Error{
		Id:   uuid.NewString(),
		Code: 404,
		Err:  errors.New("user not found"),
	}
}

func USER_ALREADY_EXISTS_ERROR() *Error {
	return &Error{
		Id:   uuid.NewString(),
		Code: 409,
		Err:  errors.New("already exists an user with this email"),
	}
}
