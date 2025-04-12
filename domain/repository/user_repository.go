package repository

import (
	"github.com/google/uuid"
	"github.com/sandroJayas/user-service/models"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	FindByID(id uuid.UUID, user *models.User) error
	Save(user *models.User) error
	SoftDelete(id uuid.UUID) error
}
