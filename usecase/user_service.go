package usecase

import (
	"github.com/google/uuid"
	"github.com/sandroJayas/user-service/domain/repository"
	"github.com/sandroJayas/user-service/models"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo}
}

func (s *UserService) Register(user *models.User) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return err
	}
	user.Password = string(hashed)
	return s.repo.CreateUser(user)
}

func (s *UserService) Login(email, password string) (*models.User, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetUserByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := s.repo.FindByID(id, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) UpdateUser(id uuid.UUID, data *models.User) (*models.User, error) {
	var user models.User
	if err := s.repo.FindByID(id, &user); err != nil {
		return nil, err
	}

	user.FirstName = data.FirstName
	user.LastName = data.LastName
	user.AddressLine1 = data.AddressLine1
	user.AddressLine2 = data.AddressLine2
	user.City = data.City
	user.PostalCode = data.PostalCode
	user.Country = data.Country
	user.PhoneNumber = data.PhoneNumber
	user.PaymentMethodID = data.PaymentMethodID

	if err := s.repo.Save(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) DeleteUser(id uuid.UUID) error {
	return s.repo.SoftDelete(id)
}
