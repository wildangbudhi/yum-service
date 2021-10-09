package usecase

import (
	"fmt"
	"log"

	"github.com/wildangbudhi/yum-service/domain"
	"github.com/wildangbudhi/yum-service/domain/v1/auth"
)

func (usecase *authUsecase) ValidateAccessToken(token *string) (*auth.ValidateAuthTokenResponse, error, domain.HTTPStatusCode) {

	if token == nil {
		return nil, fmt.Errorf("Unaothorized"), 401
	}

	var authToken auth.AuthToken = NewAuthToken(usecase.serverConfig.SecretKey, usecase.sessionRepository)

	var err error
	var tokenJWT *domain.JWT

	tokenJWT, err = domain.NewJWT(*token, usecase.serverConfig.SecretKey)

	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("Unaothorized"), 401
	}

	var resp *auth.ValidateAuthTokenResponse

	resp, err = authToken.ValidateToken(tokenJWT, false)

	if err != nil {
		return nil, fmt.Errorf("Unauthorized"), 403
	}

	return resp, nil, 200

}
