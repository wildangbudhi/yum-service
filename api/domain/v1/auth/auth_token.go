package auth

import "github.com/wildangbudhi/yum-service/domain"

type ValidateAuthTokenResponse struct {
	Role            string `json:"role"`
	IsPhoneVerified bool   `json:"is_phone_verified"`
}

type AuthToken interface {
	GenerateAuthToken(userID int, role string, isEmailVerified bool, isPhoneVerified bool) (*domain.JWT, *domain.JWT, string, error)
	ValidateToken(token *domain.JWT, isRefreshToken bool) (*ValidateAuthTokenResponse, error)
	RegenerateAuthToken(refreshToken *domain.JWT) (*domain.JWT, string, error)
	RemoveAuthToken(token *domain.JWT) error
}
