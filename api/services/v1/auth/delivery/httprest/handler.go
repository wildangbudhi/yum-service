package httprest

import (
	"github.com/gin-gonic/gin"
	"github.com/wildangbudhi/yum-service/domain/v1/auth"
)

type AuthHTTPRestHandler struct {
	authUsecase auth.AuthUsecase
}

func NewAuthHTTPRestHandler(router *gin.RouterGroup, authUsecase auth.AuthUsecase) {

	handler := AuthHTTPRestHandler{
		authUsecase: authUsecase,
	}

	router.POST("/register/customer", handler.RegisterCustomer)

}
