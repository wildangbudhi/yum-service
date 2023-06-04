package auth

import "github.com/wildangbudhi/yum-service/domain"

type AuthUsecase interface {
	RegisterCustomer(name, phoneNumber, password, apn_key, fcm_key *string) (*Customer, bool, string, string, error, domain.HTTPStatusCode)
	RegisterResto(name, phoneNumber, password, apn_key, fcm_key *string) (*Resto, bool, string, string, error, domain.HTTPStatusCode)
	AuthenticateCustomer(phoneNumber, password, apn_key, fcm_key *string) (*Customer, bool, string, string, error, domain.HTTPStatusCode)
	AuthenticateResto(phoneNumber, password, apn_key, fcm_key *string) (*Resto, bool, string, string, error, domain.HTTPStatusCode)
	ValidateAccessToken(token *string) (*ValidateAuthTokenResponse, error, domain.HTTPStatusCode)
	ResendOTPCustomer(authHeader *ValidateAuthTokenResponse) (error, domain.HTTPStatusCode)
	ResendOTPResto(authHeader *ValidateAuthTokenResponse) (error, domain.HTTPStatusCode)
	ValidateOTPCustomer(authHeader *ValidateAuthTokenResponse, otpCode *string) (*Customer, bool, string, string, error, domain.HTTPStatusCode)
	ValidateOTPResto(authHeader *ValidateAuthTokenResponse, otpCode *string) (*Resto, bool, string, string, error, domain.HTTPStatusCode)
}
