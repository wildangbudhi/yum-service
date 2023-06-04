package usecase

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func (usecase *authUsecase) HashPassword(password string) (string, error) {

	var err error
	var hashedPassword []byte

	hashedPassword, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err

}

func (usecase *authUsecase) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	if err != nil {
		log.Println(err)
	}

	return err == nil
}
