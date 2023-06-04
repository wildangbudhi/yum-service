package usecase

import (
	"fmt"

	"github.com/wildangbudhi/yum-service/domain"
	"github.com/wildangbudhi/yum-service/domain/v1/auth"
)

func (usecase *authUsecase) ResendOTPResto(authHeader *auth.ValidateAuthTokenResponse) (error, domain.HTTPStatusCode) {

	if authHeader.IsPhoneVerified {
		return fmt.Errorf("Phone number has been verified"), 400
	}

	var err error
	var resto *auth.Resto
	var otpType int = 2

	resto, err, _ = usecase.restoRepository.GetRestoByID(authHeader.UserID)

	if err != nil {
		return err, 500
	}

	var countOTP int

	countOTP, err, _ = usecase.otpRepository.CountOTPWithin30Second(authHeader.UserID.GetValue(), otpType, *resto.PhoneNumber)

	if err != nil {
		return err, 500
	}

	if countOTP > 0 {
		return fmt.Errorf("Please wait 30 second for resent otp"), 400
	}

	var phoneVerificationSID, phoneCreateVerificationResp string

	phoneVerificationSID, phoneCreateVerificationResp, err = usecase.phoneVerificationRepository.CreateAndSendOTPVerification(*resto.PhoneNumber)

	if err != nil {
		return err, 500
	}

	var newOTPLog *auth.OTP = &auth.OTP{
		ID:                         authHeader.UserID,
		Type:                       &otpType,
		PhoneNumber:                resto.PhoneNumber,
		SID:                        &phoneVerificationSID,
		CreateVerificationRespJSON: &phoneCreateVerificationResp,
	}

	err, _ = usecase.otpRepository.CreateNewOTP(newOTPLog)

	if err != nil {
		return err, 500
	}

	return nil, 200

}
