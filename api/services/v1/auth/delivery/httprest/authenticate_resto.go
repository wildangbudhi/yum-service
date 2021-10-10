package httprest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wildangbudhi/yum-service/domain"
	"github.com/wildangbudhi/yum-service/domain/v1/auth"
)

type authenticateRestoRequestBody struct {
	PhoneNumber *string `json:"phone_number" binding:"required"`
	Password    *string `json:"password" binding:"required"`
	APNKey      *string `json:"apn_key"`
	FCMKey      *string `json:"fcm_key"`
}

type authenticateRestoResponseBody struct {
	Profile               *auth.Resto `json:"profile"`
	IsPhoneNumberVerified *bool       `json:"is_phone_number_verified"`
	AccessToken           *string     `json:"access_token"`
	RefreshToken          *string     `json:"refresh_token"`
}

func (handler *AuthHTTPRestHandler) AuthenticateResto(ctx *gin.Context) {

	var err error
	var statusCode domain.HTTPStatusCode

	ctx.Header("Content-Type", "application/json")

	requestBodyData := &authenticateRestoRequestBody{}

	err = ctx.BindJSON(requestBodyData)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, domain.HTTPRestReponseBase{StatusCode: http.StatusBadRequest, Message: err.Error()})
		return
	}

	var resto *auth.Resto
	var aksesToken, refreshToken string
	var isPhoneNumberVerified bool

	resto, isPhoneNumberVerified, aksesToken, refreshToken, err, statusCode = handler.authUsecase.AuthenticateResto(
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
			Data: authenticateRestoResponseBody{
				Profile:               resto,
				IsPhoneNumberVerified: &isPhoneNumberVerified,
				AccessToken:           &aksesToken,
				RefreshToken:          &refreshToken,
			},
		},
	)

}
