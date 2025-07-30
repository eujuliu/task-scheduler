package services

import (
	"fmt"
	"log/slog"
	"scheduler/internal/entities"
	"scheduler/internal/errors"
	repos "scheduler/internal/repositories"
)

type UserService struct {
	repository repos.IUserRepository
}

func NewUserService(repo repos.IUserRepository) *UserService {
	return &UserService{
		repository: repo,
	}
}

func (s *UserService) CreateUser(username string, email string, password string) (*entities.User, error) {
	exists, _ := s.repository.GetByEmail(email)

	if exists != nil {
		return nil, errors.USER_ALREADY_EXISTS_ERROR()
	}

	user, err := entities.NewUser(username, email, password)

	if err != nil {
		slog.Debug(fmt.Sprintf("user entity creation error: %v", err))
		return nil, err
	}

	created, err := s.repository.Create(user)

	if !created {
		slog.Debug(fmt.Sprintf("user repo creation error: %v", err))
		return nil, errors.INTERNAL_SERVER_ERROR()
	}

	return user, nil
}

func (s *UserService) GetUser(email string, password string) (*entities.User, error) {
	user, err := s.repository.GetByEmail(email)

	if err != nil {
		slog.Debug(fmt.Sprintf("user find error: %v", err))
		return nil, err
	}

	ok := user.CheckPasswordHash(password)

	if !ok {
		slog.Debug(fmt.Sprintf("wrong password %s, for user %s", password, user.Id))
		return nil, errors.USER_NOT_FOUND_ERROR()
	}

	return user, nil
}
