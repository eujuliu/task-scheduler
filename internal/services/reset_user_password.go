package services

import (
	"fmt"
	"log/slog"
	"scheduler/internal/errors"
	"time"

	repos "scheduler/internal/repositories"
)

type ResetUserPasswordService struct {
	userRepository     repos.IUserRepository
	passwordRepository repos.IPasswordRecoveryRepository
}

func NewResetUserPasswordService(
	userRepo repos.IUserRepository,
	passwordRepo repos.IPasswordRecoveryRepository,
) *ResetUserPasswordService {
	return &ResetUserPasswordService{
		userRepository:     userRepo,
		passwordRepository: passwordRepo,
	}
}

func (s *ResetUserPasswordService) Execute(
	tokenId string,
	newPassword string,
) error {
	slog.Info("reset user password service started...")
	slog.Debug(fmt.Sprint("input", tokenId))

	recovery, err := s.passwordRepository.GetFirstById(tokenId)
	if err != nil {
		slog.Error(fmt.Sprintf("recovery not found error %v", err))
		return errors.RECOVERY_TOKEN_NOT_FOUND()
	}

	if !recovery.InTime(time.Now()) {
		slog.Error(fmt.Sprintf("recovery key expired error %v", err))
		return errors.RECOVERY_TOKEN_EXPIRED()
	}

	user, err := s.userRepository.GetFirstById(recovery.GetUserId())
	if err != nil {
		slog.Error(fmt.Sprintf("user not found error %v", err))
		return errors.USER_NOT_FOUND_ERROR()
	}

	err = user.SetPassword(newPassword)
	if err != nil {
		slog.Error(fmt.Sprintf("update user password error %v", err))
		return err
	}

	err = s.userRepository.Update(user)
	if err != nil {
		slog.Error(fmt.Sprintf("update user password error %v", err))
		return err
	}

	err = s.passwordRepository.Delete(recovery.GetId())
	if err != nil {
		slog.Error(fmt.Sprintf("delete recovery error %v", err))
		return err
	}

	slog.Info("reset user password service finished...")

	return nil
}
