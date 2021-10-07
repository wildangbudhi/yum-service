package auth

import "github.com/wildangbudhi/yum-service/domain"

type Customer struct {
	ID          *domain.UUID `json:"id"`
	Name        *string      `json:"name"`
	PhoneNumber *string      `json:"phone_number"`
	Password    *string      `json:"-"`
	APNKey      *string      `json:"-"`
	FCMKey      *string      `json:"-"`
}
