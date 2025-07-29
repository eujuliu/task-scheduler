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

func PASSWORD_HASHING() error {
	return &Error{
		Id:   uuid.NewString(),
		Code: 500,
		Err:  errors.New("not possible to hash password"),
	}
}

func USERNAME_INVALID() error {
	return &Error{
		Id:   uuid.NewString(),
		Code: 400,
		Err:  errors.New("username is invalid. please check the requirements"),
	}
}

func EMAIL_INVALID() error {
	return &Error{
		Id:   uuid.NewString(),
		Code: 400,
		Err:  errors.New("email is invalid. please check the requirements"),
	}
}

func PASSWORD_INVALID() error {
	return &Error{
		Id:   uuid.NewString(),
		Code: 400,
		Err:  errors.New("password is invalid. please check the requirements"),
	}
}
