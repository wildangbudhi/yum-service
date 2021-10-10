package usecase

import (
	"fmt"

	"github.com/wildangbudhi/yum-service/domain"
	"github.com/wildangbudhi/yum-service/domain/v1/auth"
)

func (usecase *authUsecase) ValidateOTP(authHeader *auth.ValidateAuthTokenResponse, otpCode *string) (error, domain.HTTPStatusCode) {

	if authHeader.IsPhoneVerified {
		return fmt.Errorf("Phone number has been verified"), 400
	}

	var userUUID *domain.UUID
	var phoneNumber string
	var otpType int
	var err error

	if authHeader.Role == "customer" {
		otpType = 1

		var customer *auth.Customer

		customer, err, _ = usecase.customerRepository.GetCustomerByID(authHeader.UserID)

		if err != nil {
			return err, 500
		}

		userUUID = customer.ID
		phoneNumber = *customer.PhoneNumber

	} else if authHeader.Role == "resto" {
		otpType = 2
	}

	var isPhoneVerified bool
	var verificationCheckRespJson, sid string

	isPhoneVerified, sid, verificationCheckRespJson, err = usecase.phoneVerificationRepository.VerifyPhone(phoneNumber, *otpCode)

	if err != nil {
		return err, 500
	}

	if !isPhoneVerified {
		return fmt.Errorf("OTP Invalid"), 400
	}

	var otpLog *auth.OTP

	otpLog, err, _ = usecase.otpRepository.GetOTP(userUUID, otpType, sid, phoneNumber)

	if err != nil {
		return err, 500
	}

	otpLog.VerificationCheckRespJSON = &verificationCheckRespJson

	err, _ = usecase.otpRepository.UpdateOTP(otpLog)

	if err != nil {
		return err, 500
	}

	var nowTimestamp *domain.Timestamp

	nowTimestamp, err = domain.NewNowTimestamp()

	if err != nil {
		return err, 500
	}

	if authHeader.Role == "customer" {

		var customer *auth.Customer

		customer, err, _ = usecase.customerRepository.GetCustomerByID(authHeader.UserID)

		if err != nil {
			return err, 500
		}

		customer.PhoneVerifiedAt = nowTimestamp

		err, _ = usecase.customerRepository.UpdateCustomer(customer)

		if err != nil {
			return err, 500
		}

	} else if authHeader.Role == "resto" {
		// update
	}

	return nil, 200

}
