package auth

import "github.com/wildangbudhi/yum-service/domain"

type AuthUsecase interface {
	RegisterCustomer(name, phoneNumber, password, apn_key, fcm_key *string) (*Customer, bool, string, string, error, domain.HTTPStatusCode)
}
