package in_memory_repos

import (
	"scheduler/internal/entities"
	"scheduler/internal/errors"
	"slices"
)

type InMemoryPasswordRecoveryRepository struct {
	tokens []entities.PasswordRecovery
}

func NewInMemoryPasswordRepository() *InMemoryPasswordRecoveryRepository {
	return &InMemoryPasswordRecoveryRepository{
		tokens: []entities.PasswordRecovery{},
	}
}

func (r *InMemoryPasswordRecoveryRepository) Get() []entities.PasswordRecovery {
	return r.tokens
}

func (r *InMemoryPasswordRecoveryRepository) GetFirstById(
	id string,
) (*entities.PasswordRecovery, error) {
	for _, token := range r.tokens {
		if token.GetId() == id {
			return &token, nil
		}
	}

	return nil, errors.RECOVERY_TOKEN_NOT_FOUND()
}

func (r *InMemoryPasswordRecoveryRepository) GetFirstByUserId(
	userId string,
) (*entities.PasswordRecovery, error) {
	for _, token := range r.tokens {
		if token.GetUserId() == userId {
			return &token, nil
		}
	}

	return nil, errors.RECOVERY_TOKEN_NOT_FOUND()
}

func (r *InMemoryPasswordRecoveryRepository) Create(
	token *entities.PasswordRecovery,
) error {
	r.tokens = append(r.tokens, *token)

	return nil
}

func (r *InMemoryPasswordRecoveryRepository) Delete(id string) error {
	index := slices.IndexFunc(r.tokens, func(u entities.PasswordRecovery) bool {
		return u.GetId() == id
	})

	r.tokens = slices.Delete(r.tokens, index, index+1)

	return nil
}
