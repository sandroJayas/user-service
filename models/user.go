package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Email    string    `json:"email" gorm:"not null"`
	Password string    `json:"password"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	FirstName string `json:"first_name" gorm:"not null"`
	LastName  string `json:"last_name" gorm:"not null"`

	AddressLine1 string `json:"address_line_1" gorm:"not null"`
	AddressLine2 string `json:"address_line_2"`
	City         string `json:"city" gorm:"not null"`
	PostalCode   string `json:"postal_code" gorm:"not null"`
	Country      string `json:"country" gorm:"not null"`
	PhoneNumber  string `json:"phone_number" gorm:"not null"`

	PaymentMethodID string `json:"payment_method_id"`
	IsDeleted       bool   `json:"is_deleted" gorm:"default:false"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}
