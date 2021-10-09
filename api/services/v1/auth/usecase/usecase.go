package usecase

import (
	"github.com/wildangbudhi/yum-service/domain/v1/auth"
	"github.com/wildangbudhi/yum-service/utils"
)

type authUsecase struct {
	serverConfig                *utils.Config
	sessionRepository           auth.SessionRepository
	phoneVerificationRepository auth.PhoneVerificationRepository
	customerRepository          auth.CustomerRepository
	otpRepository               auth.OTPRepository
}

func NewAuthUsecase(
	serverConfig *utils.Config,
	sessionRepository auth.SessionRepository,
	phoneVerificationRepository auth.PhoneVerificationRepository,
	customerRepository auth.CustomerRepository,
	otpRepository auth.OTPRepository,
) auth.AuthUsecase {

	return &authUsecase{
		serverConfig:                serverConfig,
		sessionRepository:           sessionRepository,
		phoneVerificationRepository: phoneVerificationRepository,
		customerRepository:          customerRepository,
		otpRepository:               otpRepository,
	}

}
