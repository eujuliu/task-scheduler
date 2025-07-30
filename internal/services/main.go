package services

import "scheduler/internal/entities"

type IUserService interface {
	CreateUser(username string, email string, password string) (*entities.User, error)
	GetUser(email string, password string) (*entities.User, error)
}
