package auth

import "github.com/gin-gonic/gin"

type AuthMiddlewareDelivery interface {
	ValidateAuthToken(allowedRole []string, isPhoneMustVerified bool) func(ctx *gin.Context)
}
