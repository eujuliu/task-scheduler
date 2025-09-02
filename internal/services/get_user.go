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

func (s *GetUserService) Execute(
	email string,
	password string,
) (*entities.User, error) {
	slog.Info("Get user service started...")

	user, err := s.userRepository.GetFirstByEmail(email)
	if err != nil {
		slog.Error(fmt.Sprintf("user find error: %v", err))
		return nil, err
	}

	ok := user.CheckPasswordHash(password)

	if !ok {
		slog.Error(
			fmt.Sprintf(
				"wrong password %s, for user %s",
				password,
				user.GetId(),
			),
		)
		return nil, errors.WRONG_LOGIN_DATA_ERROR()
	}

	slog.Debug(fmt.Sprintf("returned user: %+v", user))

	slog.Info("Get user service finished")
	return user, nil
}
