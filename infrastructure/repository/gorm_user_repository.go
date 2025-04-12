package repository

import (
	"github.com/google/uuid"
	"github.com/sandroJayas/user-service/models"
	"gorm.io/gorm"
)

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db}
}

func (r *GormUserRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *GormUserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ? AND is_deleted = false", email).First(&user).Error
	return &user, err
}

func (r *GormUserRepository) FindByID(id uuid.UUID, user *models.User) error {
	return r.db.First(user, "id = ? AND is_deleted = false", id).Error
}

func (r *GormUserRepository) Save(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *GormUserRepository) SoftDelete(id uuid.UUID) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).Update("is_deleted", true).Error
}
