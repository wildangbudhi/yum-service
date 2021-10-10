package httprest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wildangbudhi/yum-service/domain"
	"github.com/wildangbudhi/yum-service/domain/v1/auth"
)

type authenticateCustomerRequestBody struct {
	PhoneNumber *string `json:"phone_number" binding:"required"`
	Password    *string `json:"password" binding:"required"`
	APNKey      *string `json:"apn_key"`
	FCMKey      *string `json:"fcm_key"`
}

type authenticateCustomerResponseBody struct {
	Profile               *auth.Customer `json:"profile"`
	IsPhoneNumberVerified *bool          `json:"is_phone_number_verified"`
	AccessToken           *string        `json:"access_token"`
	RefreshToken          *string        `json:"refresh_token"`
}

func (handler *AuthHTTPRestHandler) AuthenticateCustomer(ctx *gin.Context) {

	var err error
	var statusCode domain.HTTPStatusCode

	ctx.Header("Content-Type", "application/json")

	requestBodyData := &authenticateCustomerRequestBody{}

	err = ctx.BindJSON(requestBodyData)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, domain.HTTPRestReponseBase{StatusCode: http.StatusBadRequest, Message: err.Error()})
		return
	}

	var customerData *auth.Customer
	var aksesToken, refreshToken string
	var isPhoneNumberVerified bool

	customerData, isPhoneNumberVerified, aksesToken, refreshToken, err, statusCode = handler.authUsecase.AuthenticateCustomer(
		requestBodyData.PhoneNumber,
		requestBodyData.Password,
		requestBodyData.APNKey,
		requestBodyData.FCMKey,
	)

	if err != nil {
		ctx.JSON(int(statusCode), domain.HTTPRestReponseBase{StatusCode: int(statusCode), Message: err.Error()})
		return
	}

	ctx.JSON(
		int(statusCode),
		domain.HTTPRestReponseBase{
			StatusCode: int(statusCode),
			Message:    "Success",
			Data: authenticateCustomerResponseBody{
				Profile:               customerData,
				IsPhoneNumberVerified: &isPhoneNumberVerified,
				AccessToken:           &aksesToken,
				RefreshToken:          &refreshToken,
			},
		},
	)

}
