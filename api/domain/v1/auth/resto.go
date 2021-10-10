package auth

import "github.com/wildangbudhi/yum-service/domain"

type Resto struct {
	ID              *domain.UUID      `json:"id"`
	Name            *string           `json:"name"`
	PhoneNumber     *string           `json:"phone_number"`
	Password        *string           `json:"-"`
	APNKey          *string           `json:"-"`
	FCMKey          *string           `json:"-"`
	PhoneVerifiedAt *domain.Timestamp `json:"phone_verified_at"`
}

type RestoRepository interface {
	GetRestoByID(id *domain.UUID) (*Resto, error, domain.RepositoryErrorType)
	GetRestoByPhoneNumber(phoneNumber string) (*Resto, error, domain.RepositoryErrorType)
	CreateResto(resto *Resto) (*domain.UUID, error, domain.RepositoryErrorType)
	UpdateResto(resto *Resto) (error, domain.RepositoryErrorType)
}
