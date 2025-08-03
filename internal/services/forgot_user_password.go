package services

import (
	"fmt"
	"log/slog"
	"scheduler/internal/entities"
	repos "scheduler/internal/repositories"
	"time"
)

type ForgotUserPasswordService struct {
	userRepository             repos.IUserRepository
	passwordRecoveryRepository repos.IPasswordRecoveryRepository
}

func NewForgotUserPasswordService(userRepo repos.IUserRepository, passwordRepo repos.IPasswordRecoveryRepository) *ForgotUserPasswordService {
	return &ForgotUserPasswordService{
		userRepository:             userRepo,
		passwordRecoveryRepository: passwordRepo,
	}
}

func (s *ForgotUserPasswordService) Execute(email string) (*entities.PasswordRecovery, error) {
	user, err := s.userRepository.GetByEmail(email)

	if err != nil {
		slog.Debug(fmt.Sprintf("user not found error %v", err))
		return nil, err
	}

	var recovery *entities.PasswordRecovery

	recovery, _ = s.passwordRecoveryRepository.GetByUserId(user.Id)

	if recovery != nil && recovery.ExpirationTime >= 2*time.Minute {
		slog.Debug(fmt.Sprintf("user already have an recovery token %v", recovery))
		return recovery, nil
	}

	recovery, err = entities.NewPasswordRecovery(user.Id, 5*time.Minute)

	if err != nil {
		slog.Debug(fmt.Sprintf("recovery creation error %v", err))
		return nil, err
	}

	err = s.passwordRecoveryRepository.Create(recovery)

	if err != nil {
		slog.Debug(fmt.Sprintf("recovery repository create error %v", err))
		return nil, err
	}

	return recovery, nil
}
