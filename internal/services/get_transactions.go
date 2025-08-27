package services

import (
	"scheduler/internal/entities"
	repos "scheduler/internal/repositories"
)

type GetTransactionsService struct {
	userRepository        repos.IUserRepository
	transactionRepository repos.ITransactionRepository
}

func NewGetTransactionsService(
	userRepository repos.IUserRepository,
	transactionsRepository repos.ITransactionRepository,
) *GetTransactionsService {
	return &GetTransactionsService{
		userRepository:        userRepository,
		transactionRepository: transactionsRepository,
	}
}

func (s *GetTransactionsService) Execute(userId string) []entities.Transaction {
	user, _ := s.userRepository.GetFirstById(userId)

	if user == nil {
		return []entities.Transaction{}
	}

	transactions := s.transactionRepository.GetByUserId(user.GetId())

	return transactions
}
