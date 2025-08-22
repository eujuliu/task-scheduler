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
	return fmt.Sprintf("status %d: %v", e.Code, e.Err)
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

func RECOVERY_TOKEN_NOT_FOUND() *Error {
	return &Error{
		Id:   uuid.NewString(),
		Code: http.StatusNotFound,
		Err:  errors.New("recovery token not found"),
	}
}

func RECOVERY_TOKEN_EXPIRED() *Error {
	return &Error{
		Id:   uuid.NewString(),
		Code: http.StatusUnauthorized,
		Err:  errors.New("this recovery token was expired"),
	}
}

func FINISHED_OPERATION_ERROR() *Error {
	return &Error{
		Id:   uuid.NewString(),
		Code: http.StatusUnauthorized,
		Err:  errors.New("this operation was finished, you can't change it"),
	}
}

func INVALID_FIELD_VALUE(field string) *Error {
	return &Error{
		Id:   uuid.NewString(),
		Code: http.StatusBadRequest,
		Err:  fmt.Errorf("this %s is not valid, please check the requirements", field),
	}
}

func TRANSACTION_NOT_FOUND() *Error {
	return &Error{
		Id:   uuid.NewString(),
		Code: http.StatusNotFound,
		Err:  errors.New("transaction not found"),
	}
}

func TRANSACTION_ALREADY_EXISTS_ERROR() *Error {
	return &Error{
		Id:   uuid.NewString(),
		Code: http.StatusNotFound,
		Err:  errors.New("transaction not found"),
	}
}

func NOT_ENOUGH_CREDITS_ERROR() *Error {
	return &Error{
		Id:   uuid.NewString(),
		Code: http.StatusPaymentRequired,
		Err:  errors.New("not enough credits"),
	}
}

func INVALID_TYPE() *Error {
	return &Error{
		Id:   uuid.NewString(),
		Code: http.StatusBadRequest,
		Err:  errors.New("invalid type"),
	}
}

func MISSING_PARAM_ERROR(param string) *Error {
	return &Error{
		Id:   uuid.NewString(),
		Code: http.StatusBadRequest,
		Err:  fmt.Errorf("the param %s is required for this operation", param),
	}
}

func MAX_VALUE_REACHED_ERROR() *Error {
	return &Error{
		Id:   uuid.NewString(),
		Code: http.StatusConflict,
		Err:  errors.New("max value reached"),
	}
}
