package errors

import (
	"errors"
	"fmt"
	"net/http"

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
		Code: http.StatusInternalServerError,
		Err:  errors.New("internal server error"),
	}
}

func PASSWORD_HASHING() *Error {
	return &Error{
		Id:   uuid.NewString(),
		Code: http.StatusInternalServerError,
		Err:  errors.New("not possible to hash password"),
	}
}

func USERNAME_INVALID() *Error {
	return &Error{
		Id:   uuid.NewString(),
		Code: http.StatusBadRequest,
		Err:  errors.New("username is invalid. please check the requirements"),
	}
}

func EMAIL_INVALID() *Error {
	return &Error{
		Id:   uuid.NewString(),
		Code: http.StatusBadRequest,
		Err:  errors.New("email is invalid. please check the requirements"),
	}
}

func PASSWORD_INVALID() *Error {
	return &Error{
		Id:   uuid.NewString(),
		Code: http.StatusBadRequest,
		Err:  errors.New("password is invalid. please check the requirements"),
	}
}

func USER_NOT_FOUND_ERROR() *Error {
	return &Error{
		Id:   uuid.NewString(),
		Code: http.StatusNotFound,
		Err:  errors.New("user not found"),
	}
}

func USER_ALREADY_EXISTS_ERROR() *Error {
	return &Error{
		Id:   uuid.NewString(),
		Code: http.StatusConflict,
		Err:  errors.New("already exists an user with this email"),
	}
}

func WRONG_LOGIN_DATA_ERROR() *Error {
	return &Error{
		Id:   uuid.NewString(),
		Code: http.StatusUnauthorized,
		Err:  errors.New("invalid email or password"),
	}
}
