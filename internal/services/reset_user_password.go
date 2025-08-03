package services

import (
	"fmt"
	"log/slog"
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
		slog.Debug(fmt.Sprintf("recovery not found error %v", err))
		return errors.RECOVERY_TOKEN_NOT_FOUND()
	}

	if !recovery.InTime(time.Now()) {
		slog.Debug(fmt.Sprintf("recovery key expired error %v", err))
		return errors.RECOVERY_TOKEN_EXPIRED()
	}

	user, err := s.userRepository.GetById(recovery.UserId)

	if err != nil {
		slog.Debug(fmt.Sprintf("user not found error %v", err))
		return errors.USER_NOT_FOUND_ERROR()
	}

	err = user.SetPassword(newPassword)

	if err != nil {
		slog.Debug(fmt.Sprintf("update user password error %v", err))
		return err
	}

	err = s.userRepository.Update(user)

	if err != nil {
		slog.Debug(fmt.Sprintf("update user password error %v", err))
		return err
	}

	err = s.passwordRepository.Delete(recovery.Id)

	if err != nil {
		slog.Debug(fmt.Sprintf("delete recovery error %v", err))
		return err
	}

	return nil
}
