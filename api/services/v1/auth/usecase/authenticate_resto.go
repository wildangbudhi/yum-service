package usecase

import (
	"fmt"

	"github.com/wildangbudhi/yum-service/domain"
	"github.com/wildangbudhi/yum-service/domain/v1/auth"
)

func (usecase *authUsecase) AuthenticateResto(phoneNumber, password, apn_key, fcm_key *string) (*auth.Resto, bool, string, string, error, domain.HTTPStatusCode) {

	if phoneNumber == nil {
		return nil, false, "", "", fmt.Errorf("Phone number cannot be empty"), 400
	}

	if password == nil {
		return nil, false, "", "", fmt.Errorf("Password cannot be empty"), 400
	}

	var err error
	var repositoryErrorType domain.RepositoryErrorType
	var resto *auth.Resto

	resto, err, repositoryErrorType = usecase.restoRepository.GetRestoByPhoneNumber(*phoneNumber)

	if repositoryErrorType == domain.RepositoryDataNotFound {
		return nil, false, "", "", fmt.Errorf("You haven't registered yet"), 400
	}

	if err != nil {
		return nil, false, "", "", err, 500
	}

	var passwordValid bool = usecase.CheckPasswordHash(*password, *resto.Password)

	if !passwordValid {
		return nil, false, "", "", fmt.Errorf("Phone number and password doesn't match"), 400
	}

	resto.APNKey = apn_key
	resto.FCMKey = fcm_key

	err, _ = usecase.restoRepository.UpdateResto(resto)

	if err != nil {
		return nil, false, "", "", err, 500
	}

	var authToken auth.AuthToken = NewAuthToken(usecase.serverConfig.SecretKey, usecase.sessionRepository)

	var isPhoneNumberVerified bool = false

	if resto.PhoneVerifiedAt != nil {
		isPhoneNumberVerified = true
	}

	var token, refreshToken *domain.JWT

	token, refreshToken, err = authToken.GenerateAuthToken(resto.ID.GetValue(), "resto", isPhoneNumberVerified)

	if err != nil {
		return nil, false, "", "", err, 500
	}

	return resto, isPhoneNumberVerified, token.GetToken(), refreshToken.GetToken(), nil, 200

}
