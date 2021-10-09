package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/wildangbudhi/yum-service/domain"
	"github.com/wildangbudhi/yum-service/domain/v1/auth"
)

type authRequestHeader struct {
	Authorization string `header:"Authorization" json:"Authorization" binding:"required"`
}

func (handler *authMiddlewareDelivery) ValidateAuthToken(allowedRole []string, isPhoneMustVerified bool) func(ctx *gin.Context) {

	return func(ctx *gin.Context) {

		var err error
		var resp *auth.ValidateAuthTokenResponse
		var statusCode domain.HTTPStatusCode

		requestHeader := &authRequestHeader{}

		err = ctx.BindHeader(requestHeader)

		if err != nil {
			ctx.AbortWithStatusJSON(401, domain.HTTPRestReponseBase{StatusCode: 401, Message: err.Error()})
			return
		}

		resp, err, statusCode = handler.authUsecase.ValidateAccessToken(&requestHeader.Authorization)

		if err != nil || statusCode != 200 {
			ctx.AbortWithStatusJSON(int(statusCode), domain.HTTPRestReponseBase{StatusCode: int(statusCode), Message: err.Error()})
			return
		}

		if isPhoneMustVerified && !resp.IsPhoneVerified {
			ctx.AbortWithStatusJSON(403, domain.HTTPRestReponseBase{StatusCode: 4031, Message: "Please verified your phone number"})
			return
		}

		var isRoleAllowed bool = false

		for i := 0; i < len(allowedRole); i++ {
			if resp.Role == allowedRole[i] {
				isRoleAllowed = true
				break
			}
		}

		if !isRoleAllowed {
			ctx.AbortWithStatusJSON(403, domain.HTTPRestReponseBase{StatusCode: 403, Message: "You don't have permission to this resource"})
			return
		}

		ctx.Set("AUTH_HEADER", resp)
		ctx.Next()

	}

}
