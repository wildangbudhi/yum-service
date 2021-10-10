package auth

import "github.com/wildangbudhi/yum-service/domain"

type AuthUsecase interface {
	RegisterCustomer(name, phoneNumber, password, apn_key, fcm_key *string) (*Customer, bool, string, string, error, domain.HTTPStatusCode)
	AuthenticateCustomer(phoneNumber, password, apn_key, fcm_key *string) (*Customer, bool, string, string, error, domain.HTTPStatusCode)
	ValidateAccessToken(token *string) (*ValidateAuthTokenResponse, error, domain.HTTPStatusCode)
	ResendOTPCustomer(authHeader *ValidateAuthTokenResponse) (error, domain.HTTPStatusCode)
	ValidateOTP(authHeader *ValidateAuthTokenResponse, otpCode *string) (error, domain.HTTPStatusCode)
}
