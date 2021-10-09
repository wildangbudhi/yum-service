package usecase

import "golang.org/x/crypto/bcrypt"

func (usecase *authUsecase) HashPassword(password string) (string, error) {

	var err error
	var hashedPassword []byte

	hashedPassword, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err

}
