package dto

type UpdateProfileRequest struct {
	FirstName       string `json:"first_name" binding:"required"`
	LastName        string `json:"last_name" binding:"required"`
	AddressLine1    string `json:"address_line_1" binding:"required"`
	AddressLine2    string `json:"address_line_2"`
	City            string `json:"city" binding:"required"`
	PostalCode      string `json:"postal_code" binding:"required"`
	Country         string `json:"country" binding:"required"`
	PhoneNumber     string `json:"phone_number" binding:"required"`
	PaymentMethodID string `json:"payment_method_id"`
}
