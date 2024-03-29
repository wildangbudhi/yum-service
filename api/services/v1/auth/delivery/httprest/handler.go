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
	router.POST("/register/resto", handler.RegisterResto)
	router.POST("/authenticate/customer", handler.AuthenticateCustomer)
	router.POST("/authenticate/resto", handler.AuthenticateResto)
	router.GET("/otp/resend/customer", authMiddlewareDelivery.ValidateAuthToken([]string{"customer"}, false), handler.ResendOTPCustomer)
	router.GET("/otp/resend/resto", authMiddlewareDelivery.ValidateAuthToken([]string{"resto"}, false), handler.ResendOTPResto)
	router.POST("/otp/validate/customer", authMiddlewareDelivery.ValidateAuthToken([]string{"customer"}, false), handler.ValidateOTPCustomer)
	router.POST("/otp/validate/resto", authMiddlewareDelivery.ValidateAuthToken([]string{"resto"}, false), handler.ValidateOTPResto)

}
