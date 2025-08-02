package services

import (
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
		return nil, err
	}

	var recovery *entities.PasswordRecovery

	recovery, _ = s.passwordRecoveryRepository.GetByUserId(user.Id)

	if recovery != nil && recovery.ExpirationTime >= 2*time.Minute {
		return recovery, nil
	}

	recovery, err = entities.NewPasswordRecovery(user.Id, 5*time.Minute)

	if err != nil {
		return nil, err
	}

	err = s.passwordRecoveryRepository.Create(recovery)

	if err != nil {
		return nil, err
	}

	return recovery, nil
}
