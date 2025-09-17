package services

import (
	"fmt"
	"log/slog"
	"scheduler/internal/entities"
	"scheduler/internal/interfaces"
	"time"
)

type ForgotUserPasswordService struct {
	userRepository             interfaces.IUserRepository
	passwordRecoveryRepository interfaces.IPasswordRecoveryRepository
}

func NewForgotUserPasswordService(
	userRepo interfaces.IUserRepository,
	passwordRepo interfaces.IPasswordRecoveryRepository,
) *ForgotUserPasswordService {
	return &ForgotUserPasswordService{
		userRepository:             userRepo,
		passwordRecoveryRepository: passwordRepo,
	}
}

func (s *ForgotUserPasswordService) Execute(
	email string,
) (*entities.PasswordRecovery, error) {
	slog.Info("forgot password service started...")
	slog.Debug(fmt.Sprint("input ", email))

	user, err := s.userRepository.GetFirstByEmail(email)
	if err != nil {
		slog.Error(fmt.Sprintf("user not found error %v", err))
		return nil, err
	}

	var recovery *entities.PasswordRecovery

	recovery, _ = s.passwordRecoveryRepository.GetFirstByUserId(user.GetId())

	if recovery != nil && recovery.GetExpiration() >= 2*time.Minute {
		slog.Error(
			fmt.Sprintf("user already have an recovery token %v", recovery),
		)
		return recovery, nil
	}

	recovery, err = entities.NewPasswordRecovery(user.GetId(), 5*time.Minute)
	if err != nil {
		slog.Error(fmt.Sprintf("recovery creation error %v", err))
		return nil, err
	}

	err = s.passwordRecoveryRepository.Create(recovery)
	if err != nil {
		slog.Error(fmt.Sprintf("recovery repository create error %v", err))
		return nil, err
	}

	slog.Info("forgot password service finished...")
	slog.Debug(fmt.Sprint("input", recovery))

	return recovery, nil
}
