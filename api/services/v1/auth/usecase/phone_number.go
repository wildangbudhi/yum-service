package usecase

import (
	"fmt"
	"log"
	"regexp"
)

func (usecase *authUsecase) ValidatePhoneNumberFormat(phoneNumber string) bool {
	re := regexp.MustCompile(`\+[0-9]+`)
	return re.MatchString(phoneNumber)
}

func (usecase *authUsecase) SanitizePhoneNumber(phoneNumber *string) error {

	reg, err := regexp.Compile("[^0-9+]+")

	if err != nil {
		log.Println(err)
		return fmt.Errorf("Failed to process phone number")
	}

	*phoneNumber = reg.ReplaceAllString(*phoneNumber, "")

	if (*phoneNumber)[0] == '0' {
		*phoneNumber = (*phoneNumber)[1:]
		*phoneNumber = "+62" + *phoneNumber
	}

	return nil

}
