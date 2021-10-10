package auth

type PhoneVerificationRepository interface {
	CreateAndSendOTPVerification(phoneNumber string) (string, string, error)
	VerifyPhone(phoneNumber, otpCode string) (bool, string, string, error)
}
