package usecase

import (
	"fmt"
	"log"

	"github.com/wildangbudhi/yum-service/domain"
	"github.com/wildangbudhi/yum-service/domain/v1/auth"
)

func (usecase *authUsecase) RegisterCustomer(name, phoneNumber, password, apn_key, fcm_key *string) (*auth.Customer, bool, string, string, error, domain.HTTPStatusCode) {

	if name == nil {
		return nil, false, "", "", fmt.Errorf("Name cannot be empty"), 400
	}

	if phoneNumber == nil {
		return nil, false, "", "", fmt.Errorf("Phone number cannot be empty"), 400
	}

	if password == nil {
		return nil, false, "", "", fmt.Errorf("Password cannot be empty"), 400
	}

	var err error
	var repositoryErrorType domain.RepositoryErrorType

	err = usecase.SanitizePhoneNumber(phoneNumber)

	if err != nil {
		return nil, false, "", "", err, 400
	}

	var isPhoneNumberValid bool = usecase.ValidatePhoneNumberFormat(*phoneNumber)

	if !isPhoneNumberValid {
		return nil, false, "", "", fmt.Errorf("Please use +62xxxxxx phone number format"), 400
	}

	var hashedPassword string

	hashedPassword, err = usecase.HashPassword(*password)

	if err != nil {
		log.Println(err)
		return nil, false, "", "", fmt.Errorf("Failed to process registration, please try again"), 500
	}

	var newCustomer *auth.Customer = &auth.Customer{
		Name:        name,
		PhoneNumber: phoneNumber,
		Password:    &hashedPassword,
		APNKey:      apn_key,
		FCMKey:      fcm_key,
	}

	newCustomer.ID, err, repositoryErrorType = usecase.customerRepository.CreateCustomer(newCustomer)

	if repositoryErrorType == domain.RepositoryCreateDataFailed {
		return nil, false, "", "", err, 400
	}

	if err != nil {
		return nil, false, "", "", err, 500
	}

	var phoneVerificationSID, phoneCreateVerificationResp string

	phoneVerificationSID, phoneCreateVerificationResp, err = usecase.phoneVerificationRepository.CreateAndSendOTPVerification(*phoneNumber)

	if err != nil {
		return nil, false, "", "", err, 500
	}

	var otpType int = domain.OTPCustomertype

	var newOTPLog *auth.OTP = &auth.OTP{
		ID:                         newCustomer.ID,
		Type:                       &otpType,
		SID:                        &phoneVerificationSID,
		CreateVerificationRespJSON: &phoneCreateVerificationResp,
	}

	err, _ = usecase.otpRepository.CreateNewOTP(newOTPLog)

	if err != nil {
		return nil, false, "", "", err, 500
	}

	var authToken auth.AuthToken = NewAuthToken(usecase.serverConfig.SecretKey, usecase.sessionRepository)

	var isPhoneNumberVerified bool = false

	if newCustomer.PhoneVerifiedAt != nil {
		isPhoneNumberVerified = true
	}

	var token, refreshToken *domain.JWT

	token, refreshToken, err = authToken.GenerateAuthToken(newCustomer.ID.GetValue(), "customer", isPhoneNumberVerified)

	if err != nil {
		return nil, false, "", "", err, 500
	}

	return newCustomer, isPhoneNumberVerified, token.GetToken(), refreshToken.GetToken(), nil, 200

}
