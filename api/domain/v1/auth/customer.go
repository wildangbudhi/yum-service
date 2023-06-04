package auth

import "github.com/wildangbudhi/yum-service/domain"

type Customer struct {
	ID              *domain.UUID      `json:"id"`
	Name            *string           `json:"name"`
	PhoneNumber     *string           `json:"phone_number"`
	Password        *string           `json:"-"`
	APNKey          *string           `json:"-"`
	FCMKey          *string           `json:"-"`
	PhoneVerifiedAt *domain.Timestamp `json:"phone_verified_at"`
}

type CustomerRepository interface {
	GetCustomerByID(id *domain.UUID) (*Customer, error, domain.RepositoryErrorType)
	GetCustomerByPhoneNumber(phoneNumber string) (*Customer, error, domain.RepositoryErrorType)
	CreateCustomer(customer *Customer) (*domain.UUID, error, domain.RepositoryErrorType)
	UpdateCustomer(customer *Customer) (error, domain.RepositoryErrorType)
}
