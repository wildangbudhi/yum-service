package middleware

import "github.com/wildangbudhi/yum-service/domain/v1/auth"

type authMiddlewareDelivery struct {
	authUsecase auth.AuthUsecase
}

func NewAuthMiddlewareDelivery(authUsecase auth.AuthUsecase) auth.AuthMiddlewareDelivery {
	return &authMiddlewareDelivery{
		authUsecase: authUsecase,
	}
}
