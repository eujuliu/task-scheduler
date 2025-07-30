package services

import (
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
		return nil, err
	}

	created := s.repository.Create(user)

	if !created {
		return nil, errors.INTERNAL_SERVER_ERROR()
	}

	return user, nil
}
