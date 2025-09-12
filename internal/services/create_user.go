package services

import (
	"fmt"
	"log/slog"
	"scheduler/internal/entities"
	"scheduler/internal/errors"
	"scheduler/internal/interfaces"
)

type CreateUserService struct {
	userRepository         interfaces.IUserRepository
	customerPaymentGateway interfaces.ICustomerPaymentGateway
}

func NewCreateUserService(
	userRepo interfaces.IUserRepository,
	customerPaymentGateway interfaces.ICustomerPaymentGateway,
) *CreateUserService {
	return &CreateUserService{
		userRepository:         userRepo,
		customerPaymentGateway: customerPaymentGateway,
	}
}

func (s *CreateUserService) Execute(
	username string,
	email string,
	password string,
) (*entities.User, error) {
	slog.Info("create user service started...")
	slog.Debug(fmt.Sprint("input ", username,
		email))

	exists, _ := s.userRepository.GetFirstByEmail(email)

	if exists != nil {
		return nil, errors.USER_ALREADY_EXISTS_ERROR()
	}

	user, err := entities.NewUser(username, email, password)
	if err != nil {
		slog.Error(fmt.Sprintf("user entity creation error: %v", err))
		return nil, err
	}

	err = s.userRepository.Create(user)
	if err != nil {
		slog.Error(fmt.Sprintf("user repo creation error: %v", err))
		return nil, errors.INTERNAL_SERVER_ERROR()
	}

	_, err = s.customerPaymentGateway.Create(user.GetId(), user.GetUsername(), user.GetEmail(), nil)
	if err != nil {
		return nil, err
	}

	slog.Debug(fmt.Sprintf("returned user: %+v", user))

	slog.Info("create user service finished")

	return user, nil
}
