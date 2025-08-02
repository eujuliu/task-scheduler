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

func (p *InMemoryPasswordRecoveryRepository) Get() []entities.PasswordRecovery {
	return p.tokens
}

func (p *InMemoryPasswordRecoveryRepository) GetById(id string) (*entities.PasswordRecovery, error) {
	for _, token := range p.tokens {
		if token.Id == id {
			return &token, nil
		}
	}

	return nil, errors.RECOVERY_TOKEN_NOT_FOUND()
}

func (p *InMemoryPasswordRecoveryRepository) GetByUserId(userId string) (*entities.PasswordRecovery, error) {
	for _, token := range p.tokens {
		if token.UserId == userId {
			return &token, nil
		}
	}

	return nil, errors.RECOVERY_TOKEN_NOT_FOUND()
}

func (p *InMemoryPasswordRecoveryRepository) Create(token *entities.PasswordRecovery) error {
	p.tokens = append(p.tokens, *token)

	return nil
}

func (p *InMemoryPasswordRecoveryRepository) Delete(id string) error {
	index := slices.IndexFunc(p.tokens, func(u entities.PasswordRecovery) bool {
		return u.Id == id
	})

	p.tokens = slices.Delete(p.tokens, index, index+1)

	return nil
}
