package httprest

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wildangbudhi/yum-service/domain"
	"github.com/wildangbudhi/yum-service/domain/v1/auth"
)

type validateOTPRestoRequestBody struct {
	OTPCode *string `json:"otp_code" binding:"required"`
}

type validateOTPRestoResponseBody struct {
	Profile               *auth.Resto `json:"profile"`
	IsPhoneNumberVerified *bool       `json:"is_phone_number_verified"`
	AccessToken           *string     `json:"access_token"`
	RefreshToken          *string     `json:"refresh_token"`
}

func (handler *AuthHTTPRestHandler) ValidateOTPResto(ctx *gin.Context) {

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

	requestBodyData := &validateOTPRestoRequestBody{}

	err = ctx.BindJSON(requestBodyData)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, domain.HTTPRestReponseBase{StatusCode: http.StatusBadRequest, Message: err.Error()})
		return
	}

	var restoData *auth.Resto
	var aksesToken, refreshToken string
	var isPhoneNumberVerified bool

	restoData, isPhoneNumberVerified, aksesToken, refreshToken, err, statusCode = handler.authUsecase.ValidateOTPResto(authHeader, requestBodyData.OTPCode)

	if err != nil {
		ctx.JSON(int(statusCode), domain.HTTPRestReponseBase{StatusCode: int(statusCode), Message: err.Error()})
		return
	}

	ctx.JSON(
		int(statusCode),
		domain.HTTPRestReponseBase{
			StatusCode: int(statusCode),
			Message:    "Success",
			Data: validateOTPRestoResponseBody{
				Profile:               restoData,
				IsPhoneNumberVerified: &isPhoneNumberVerified,
				AccessToken:           &aksesToken,
				RefreshToken:          &refreshToken,
			},
		},
	)

}
