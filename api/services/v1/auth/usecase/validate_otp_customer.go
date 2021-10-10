package usecase

import (
	"fmt"

	"github.com/wildangbudhi/yum-service/domain"
	"github.com/wildangbudhi/yum-service/domain/v1/auth"
)

func (usecase *authUsecase) ValidateOTPCustomer(authHeader *auth.ValidateAuthTokenResponse, otpCode *string) (*auth.Customer, bool, string, string, error, domain.HTTPStatusCode) {

	if authHeader.IsPhoneVerified {
		return nil, false, "", "", fmt.Errorf("Phone number has been verified"), 400
	}

	var otpType int = 1
	var err error

	var customer *auth.Customer

	customer, err, _ = usecase.customerRepository.GetCustomerByID(authHeader.UserID)

	if err != nil {
		return nil, false, "", "", err, 500
	}

	var isPhoneVerified bool
	var verificationCheckRespJson, sid string

	isPhoneVerified, sid, verificationCheckRespJson, err = usecase.phoneVerificationRepository.VerifyPhone(*customer.PhoneNumber, *otpCode)

	if err != nil {
		return nil, false, "", "", err, 500
	}

	if !isPhoneVerified {
		return nil, false, "", "", fmt.Errorf("OTP Invalid"), 400
	}

	var otpLog *auth.OTP

	otpLog, err, _ = usecase.otpRepository.GetOTP(customer.ID, otpType, sid, *customer.PhoneNumber)

	if err != nil {
		return nil, false, "", "", err, 500
	}

	otpLog.VerificationCheckRespJSON = &verificationCheckRespJson

	err, _ = usecase.otpRepository.UpdateOTP(otpLog)

	if err != nil {
		return nil, false, "", "", err, 500
	}

	var nowTimestamp *domain.Timestamp

	nowTimestamp, err = domain.NewNowTimestamp()

	if err != nil {
		return nil, false, "", "", err, 500
	}

	customer.PhoneVerifiedAt = nowTimestamp

	err, _ = usecase.customerRepository.UpdateCustomer(customer)

	if err != nil {
		return nil, false, "", "", err, 500
	}

	var authToken auth.AuthToken = NewAuthToken(usecase.serverConfig.SecretKey, usecase.sessionRepository)

	var token, refreshToken *domain.JWT

	token, refreshToken, err = authToken.GenerateAuthToken(customer.ID.GetValue(), "customer", true)

	if err != nil {
		return nil, false, "", "", err, 500
	}

	return customer, true, token.GetToken(), refreshToken.GetToken(), nil, 200

}
