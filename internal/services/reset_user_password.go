package services

import (
	"scheduler/internal/errors"
	repos "scheduler/internal/repositories"
	"time"
)

type ResetUserPasswordService struct {
	userRepository     repos.IUserRepository
	passwordRepository repos.IPasswordRecoveryRepository
}

func NewResetUserPasswordService(userRepo repos.IUserRepository, passwordRepo repos.IPasswordRecoveryRepository) *ResetUserPasswordService {
	return &ResetUserPasswordService{
		userRepository:     userRepo,
		passwordRepository: passwordRepo,
	}
}

func (s *ResetUserPasswordService) Execute(tokenId string, newPassword string) error {

	recovery, err := s.passwordRepository.GetById(tokenId)

	if err != nil {
		return errors.RECOVERY_TOKEN_NOT_FOUND()
	}

	if !recovery.InTime(time.Now()) {
		return errors.RECOVERY_TOKEN_EXPIRED()
	}

	user, err := s.userRepository.GetById(recovery.UserId)

	if err != nil {
		return errors.USER_NOT_FOUND_ERROR()
	}

	err = user.SetPassword(newPassword)

	if err != nil {
		return err
	}

	err = s.userRepository.Update(user)

	if err != nil {
		return err
	}

	err = s.passwordRepository.Delete(recovery.Id)

	if err != nil {
		return err
	}

	return nil
}
