package httprest

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wildangbudhi/yum-service/domain"
	"github.com/wildangbudhi/yum-service/domain/v1/auth"
)

type validateOTPRequestBody struct {
	OTPCode *string `json:"otp_code" binding:"required"`
}

func (handler *AuthHTTPRestHandler) ValidateOTP(ctx *gin.Context) {

	var err error
	var statusCode domain.HTTPStatusCode

	ctx.Header("Content-Type", "application/json")

	var authHeaderInterface interface{}
	var isAuthHeaderExists bool = false

	authHeaderInterface, isAuthHeaderExists = ctx.Get("AUTH_HEADER")

	if !isAuthHeaderExists {
		log.Println("Auth header not found")
		ctx.JSON(http.StatusBadRequest, domain.HTTPRestReponseBase{StatusCode: http.StatusBadRequest, Message: "Unauthorized"})
		return
	}

	var isConversionOK bool = false
	var authHeader *auth.ValidateAuthTokenResponse

	authHeader, isConversionOK = authHeaderInterface.(*auth.ValidateAuthTokenResponse)

	if !isConversionOK {
		log.Println("Cannot convert interface{} to *auth.ValidateAuthTokenResponse")
		ctx.JSON(http.StatusBadRequest, domain.HTTPRestReponseBase{StatusCode: http.StatusBadRequest, Message: "Unauthorized"})
		return
	}

	requestBodyData := &validateOTPRequestBody{}

	err = ctx.BindJSON(requestBodyData)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, domain.HTTPRestReponseBase{StatusCode: http.StatusBadRequest, Message: err.Error()})
		return
	}

	err, statusCode = handler.authUsecase.ValidateOTP(authHeader, requestBodyData.OTPCode)

	if err != nil {
		ctx.JSON(int(statusCode), domain.HTTPRestReponseBase{StatusCode: int(statusCode), Message: err.Error()})
		return
	}

	ctx.JSON(int(200), domain.HTTPRestReponseBase{StatusCode: int(200), Message: "Success"})

}
