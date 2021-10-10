package usecase

import (
	"fmt"

	"github.com/wildangbudhi/yum-service/domain"
	"github.com/wildangbudhi/yum-service/domain/v1/auth"
)

func (usecase *authUsecase) ValidateOTPResto(authHeader *auth.ValidateAuthTokenResponse, otpCode *string) (*auth.Resto, bool, string, string, error, domain.HTTPStatusCode) {

	if authHeader.IsPhoneVerified {
		return nil, false, "", "", fmt.Errorf("Phone number has been verified"), 400
	}

	var otpType int = domain.OTPRestotype
	var err error

	var resto *auth.Resto

	resto, err, _ = usecase.restoRepository.GetRestoByID(authHeader.UserID)

	if err != nil {
		return nil, false, "", "", err, 500
	}

	var isPhoneVerified bool
	var verificationCheckRespJson, sid string

	isPhoneVerified, sid, verificationCheckRespJson, err = usecase.phoneVerificationRepository.VerifyPhone(*resto.PhoneNumber, *otpCode)

	if err != nil {
		return nil, false, "", "", err, 500
	}

	if !isPhoneVerified {
		return nil, false, "", "", fmt.Errorf("OTP Invalid"), 400
	}

	var otpLog *auth.OTP

	otpLog, err, _ = usecase.otpRepository.GetOTP(resto.ID, otpType, sid, *resto.PhoneNumber)

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

	resto.PhoneVerifiedAt = nowTimestamp

	err, _ = usecase.restoRepository.UpdateResto(resto)

	if err != nil {
		return nil, false, "", "", err, 500
	}

	var authToken auth.AuthToken = NewAuthToken(usecase.serverConfig.SecretKey, usecase.sessionRepository)

	var token, refreshToken *domain.JWT

	token, refreshToken, err = authToken.GenerateAuthToken(resto.ID.GetValue(), "resto", true)

	if err != nil {
		return nil, false, "", "", err, 500
	}

	return resto, true, token.GetToken(), refreshToken.GetToken(), nil, 200

}
