package sms

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/verify/v2"
	"github.com/wildangbudhi/yum-service/domain/v1/auth"
	"github.com/wildangbudhi/yum-service/utils"
)

type phoneVerificationRepository struct {
	config     utils.Config
	smsService *twilio.RestClient
}

func NewPhoneVerificationRepository(config utils.Config, smsService *twilio.RestClient) auth.PhoneVerificationRepository {
	return &phoneVerificationRepository{
		config:     config,
		smsService: smsService,
	}
}

func (repo *phoneVerificationRepository) CreateAndSendOTPVerification(phoneNumber string) (string, string, error) {

	var err error
	var channel string = "sms"

	res, err := repo.smsService.VerifyV2.CreateVerification(
		repo.config.TwilioPhoneNumberVerificationServiceID,
		&openapi.CreateVerificationParams{
			Channel: &channel,
			To:      &phoneNumber,
		},
	)

	if err != nil {
		log.Println(err)
		return "", "", fmt.Errorf("Services Unavailable")
	}

	var verificationID string = *res.Sid

	byteJSON, err := json.Marshal(res)

	if err != nil {
		log.Println(err)
		return "", "", fmt.Errorf("Services Unavailable")
	}

	return verificationID, string(byteJSON), nil

}

func (repo *phoneVerificationRepository) VerifyPhone(phoneNumber, otpCode string) (bool, string, error) {

	var err error

	res, err := repo.smsService.VerifyV2.CreateVerificationCheck(
		repo.config.TwilioPhoneNumberVerificationServiceID,
		&openapi.CreateVerificationCheckParams{
			To:   &phoneNumber,
			Code: &otpCode,
		},
	)

	byteJSON, err := json.Marshal(res)

	if err != nil {
		log.Println(err)
		return false, "", fmt.Errorf("Services Unavailable")
	}

	return *res.Valid, string(byteJSON), nil

}
