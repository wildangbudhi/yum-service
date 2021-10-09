package httprest

import (
	"github.com/gin-gonic/gin"
	"github.com/wildangbudhi/yum-service/domain/v1/auth"
)

type AuthHTTPRestHandler struct {
	authMiddlewareDelivery auth.AuthMiddlewareDelivery
	authUsecase            auth.AuthUsecase
}

func NewAuthHTTPRestHandler(router *gin.RouterGroup, authMiddlewareDelivery auth.AuthMiddlewareDelivery, authUsecase auth.AuthUsecase) {

	handler := AuthHTTPRestHandler{
		authMiddlewareDelivery: authMiddlewareDelivery,
		authUsecase:            authUsecase,
	}

	router.POST("/register/customer", handler.RegisterCustomer)
	router.POST("/authenticate/customer", handler.AuthenticateCustomer)
	router.GET("/otp/resend", authMiddlewareDelivery.ValidateAuthToken([]string{"customer", "role"}, false), handler.ResendOTP)

}
