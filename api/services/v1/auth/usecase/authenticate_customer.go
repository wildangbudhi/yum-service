package usecase

import (
	"fmt"

	"github.com/wildangbudhi/yum-service/domain"
	"github.com/wildangbudhi/yum-service/domain/v1/auth"
)

func (usecase *authUsecase) AuthenticateCustomer(phoneNumber, password, apn_key, fcm_key *string) (*auth.Customer, bool, string, string, error, domain.HTTPStatusCode) {

	if phoneNumber == nil {
		return nil, false, "", "", fmt.Errorf("Phone number cannot be empty"), 400
	}

	if password == nil {
		return nil, false, "", "", fmt.Errorf("Password cannot be empty"), 400
	}

	var err error
	var repositoryErrorType domain.RepositoryErrorType
	var customer *auth.Customer

	customer, err, repositoryErrorType = usecase.customerRepository.GetCustomerByPhoneNumber(*phoneNumber)

	if repositoryErrorType == domain.RepositoryDataNotFound {
		return nil, false, "", "", fmt.Errorf("You haven't registered yet"), 400
	}

	if err != nil {
		return nil, false, "", "", err, 500
	}

	var passwordValid bool = usecase.CheckPasswordHash(*password, *customer.Password)

	if !passwordValid {
		return nil, false, "", "", fmt.Errorf("Phone number and password doesn't match"), 400
	}

	customer.APNKey = apn_key
	customer.FCMKey = fcm_key

	err, _ = usecase.customerRepository.UpdateCustomer(customer)

	if err != nil {
		return nil, false, "", "", err, 500
	}

	var authToken auth.AuthToken = NewAuthToken(usecase.serverConfig.SecretKey, usecase.sessionRepository)

	var isPhoneNumberVerified bool = false

	if customer.PhoneVerifiedAt != nil {
		isPhoneNumberVerified = true
	}

	var token, refreshToken *domain.JWT

	token, refreshToken, err = authToken.GenerateAuthToken(customer.ID.GetValue(), "customer", isPhoneNumberVerified)

	if err != nil {
		return nil, false, "", "", err, 500
	}

	return customer, isPhoneNumberVerified, token.GetToken(), refreshToken.GetToken(), nil, 200

}
