package services

import (
	"scheduler/internal/entities"
	"scheduler/internal/interfaces"
)

type GetTransactionsService struct {
	userRepository        interfaces.IUserRepository
	transactionRepository interfaces.ITransactionRepository
}

func NewGetTransactionsService(
	userRepository interfaces.IUserRepository,
	transactionsRepository interfaces.ITransactionRepository,
) *GetTransactionsService {
	return &GetTransactionsService{
		userRepository:        userRepository,
		transactionRepository: transactionsRepository,
	}
}

func (s *GetTransactionsService) Execute(
	userId string,
	offset *int,
	limit *int,
	orderBy *string,
) []entities.Transaction {
	user, _ := s.userRepository.GetFirstById(userId)

	if user == nil {
		return []entities.Transaction{}
	}

	transactions := s.transactionRepository.GetByUserId(user.GetId(), offset, limit, orderBy)

	return transactions
}
