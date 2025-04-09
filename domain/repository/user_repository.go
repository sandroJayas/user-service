package repository

import "github.com/sandroJayas/user-service/models"

type UserRepository interface {
	CreateUser(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	FindByID(id string, user *models.User) error
	Save(user *models.User) error
	SoftDelete(id string) error
}
