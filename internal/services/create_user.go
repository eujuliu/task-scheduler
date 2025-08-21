package services

import (
	"fmt"
	"log/slog"
	"scheduler/internal/entities"
	"scheduler/internal/errors"
	repos "scheduler/internal/repositories"
)

type CreateUserService struct {
	userRepository repos.IUserRepository
}

func NewCreateUserService(userRepo repos.IUserRepository) *CreateUserService {
	return &CreateUserService{
		userRepository: userRepo,
	}
}

func (s *CreateUserService) Execute(username string, email string, password string) (*entities.User, error) {
	exists, _ := s.userRepository.GetFirstByEmail(email)

	if exists != nil {
		return nil, errors.USER_ALREADY_EXISTS_ERROR()
	}

	user, err := entities.NewUser(username, email, password)

	if err != nil {
		slog.Debug(fmt.Sprintf("user entity creation error: %v", err))
		return nil, err
	}

	err = s.userRepository.Create(user)

	if err != nil {
		slog.Debug(fmt.Sprintf("user repo creation error: %v", err))
		return nil, errors.INTERNAL_SERVER_ERROR()
	}

	return user, nil
}
