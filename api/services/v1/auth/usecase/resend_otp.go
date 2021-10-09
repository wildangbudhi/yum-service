package usecase

import (
	"fmt"

	"github.com/wildangbudhi/yum-service/domain"
	"github.com/wildangbudhi/yum-service/domain/v1/auth"
)

func (usecase *authUsecase) ResendOTP(authHeader *auth.ValidateAuthTokenResponse) (error, domain.HTTPStatusCode) {

	if authHeader.IsPhoneVerified {
		return fmt.Errorf("Phone number has been verified"), 400
	}

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

		phoneNumber = *customer.PhoneNumber

	} else if authHeader.Role == "resto" {
		otpType = 2
	}

	var countOTP int

	countOTP, err, _ = usecase.otpRepository.CountOTPWithin30Second(authHeader.UserID.GetValue(), otpType, phoneNumber)

	if err != nil {
		return err, 500
	}

	if countOTP > 0 {
		return fmt.Errorf("Please wait 30 second for resent otp"), 400
	}

	var phoneVerificationSID, phoneCreateVerificationResp string

	phoneVerificationSID, phoneCreateVerificationResp, err = usecase.phoneVerificationRepository.CreateAndSendOTPVerification(phoneNumber)

	if err != nil {
		return err, 500
	}

	var newOTPLog *auth.OTP = &auth.OTP{
		ID:                         authHeader.UserID,
		Type:                       &otpType,
		PhoneNumber:                &phoneNumber,
		SID:                        &phoneVerificationSID,
		CreateVerificationRespJSON: &phoneCreateVerificationResp,
	}

	err, _ = usecase.otpRepository.CreateNewOTP(newOTPLog)

	if err != nil {
		return err, 500
	}

	return nil, 200

}
