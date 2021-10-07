package auth

import "github.com/wildangbudhi/yum-service/domain"

type OTP struct {
	ID                         *domain.UUID
	Type                       *int
	SID                        *string
	CreateVerificationRespJSON *string
	VerificationCheckRespJSON  *string
}

type OTPRepository interface {
	CountOTPWithin30Second(id string, userType int) (int, error, domain.RepositoryErrorType)
	CreateNewOTP(otp *OTP) (error, domain.RepositoryErrorType)
	UpdateOTP(otp *OTP) (error, domain.RepositoryErrorType)
}
