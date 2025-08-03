package services

import (
	"fmt"
	"log/slog"
	"scheduler/internal/entities"
	"scheduler/internal/errors"
	repos "scheduler/internal/repositories"
)

type GetUserService struct {
	userRepository repos.IUserRepository
}

func NewGetUserService(userRepo repos.IUserRepository) *GetUserService {
	return &GetUserService{
		userRepository: userRepo,
	}
}

func (s *GetUserService) Execute(email string, password string) (*entities.User, error) {
	user, err := s.userRepository.GetByEmail(email)

	if err != nil {
		slog.Debug(fmt.Sprintf("user find error: %v", err))
		return nil, err
	}

	ok := user.CheckPasswordHash(password)

	if !ok {
		slog.Debug(fmt.Sprintf("wrong password %s, for user %s", password, user.Id))
		return nil, errors.WRONG_LOGIN_DATA_ERROR()
	}

	return user, nil
}
