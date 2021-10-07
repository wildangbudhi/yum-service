package auth

import "github.com/wildangbudhi/yum-service/domain"

type Resto struct {
	ID          *domain.UUID `json:"id"`
	Name        *string      `json:"name"`
	PhoneNumber *string      `json:"phone_number"`
	Password    *string      `json:"-"`
	APNKey      *string      `json:"-"`
	FCMKey      *string      `json:"-"`
}
