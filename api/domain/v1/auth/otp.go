package auth

import "github.com/wildangbudhi/yum-service/domain"

type OTP struct {
	ID                         *domain.UUID
	Type                       *int
	PhoneNumber                *string
	SID                        *string
	CreateVerificationRespJSON *string
	VerificationCheckRespJSON  *string
}

type OTPRepository interface {
	CountOTPWithin30Second(id string, userType int, phoneNumber string) (int, error, domain.RepositoryErrorType)
	CreateNewOTP(otp *OTP) (error, domain.RepositoryErrorType)
	GetOTP(userID *domain.UUID, otpType int, sid, phoneNumber string) (*OTP, error, domain.RepositoryErrorType)
	UpdateOTP(otp *OTP) (error, domain.RepositoryErrorType)
}
