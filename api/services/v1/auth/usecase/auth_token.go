package usecase

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/wildangbudhi/yum-service/domain"
	"github.com/wildangbudhi/yum-service/domain/v1/auth"
)

type authToken struct {
	secretKey         []byte
	sessionRepository auth.SessionRepository
}

func NewAuthToken(secretKey []byte, sessionRepository auth.SessionRepository) auth.AuthToken {
	return &authToken{
		secretKey:         secretKey,
		sessionRepository: sessionRepository,
	}
}

func (obj *authToken) generateSecretKey() (string, error) {

	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-!@#$%^&*()+=_"

	ret := make([]byte, 64)

	for i := 0; i < len(ret); i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))

		if err != nil {
			return "", err
		}

		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil

}

func (obj *authToken) GenerateAuthToken(userID int, role string, isPhoneVerified bool) (*domain.JWT, *domain.JWT, string, error) {

	var err error

	var token, refreshToken *domain.JWT
	var tokenPayload, refreshTokenPayload jwt.MapClaims
	var tokenUUID *domain.UUID = domain.NewUUID()
	var refreshTokenUUID *domain.UUID = domain.NewUUID()

	tokenPayload = jwt.MapClaims{}
	tokenPayload["user_id"] = userID
	tokenPayload["role"] = role
	tokenPayload["uuid"] = tokenUUID.GetValue()
	tokenPayload["exp"] = time.Now().Add(time.Hour * 24).Unix()

	token, err = domain.NewJWTFromPayload(tokenPayload, obj.secretKey)

	if err != nil {
		log.Println(err)
		return nil, nil, "", fmt.Errorf("Failed to Generate Session")
	}

	refreshTokenPayload = jwt.MapClaims{}
	refreshTokenPayload["user_id"] = userID
	refreshTokenPayload["role"] = role
	refreshTokenPayload["uuid"] = refreshTokenUUID.GetValue()

	refreshToken, err = domain.NewJWTFromPayload(refreshTokenPayload, obj.secretKey)

	if err != nil {
		log.Println(err)
		return nil, nil, "", fmt.Errorf("Failed to Generate Session")
	}

	token, err = domain.NewJWTFromPayload(tokenPayload, obj.secretKey)

	if err != nil {
		log.Println(err)
		return nil, nil, "", fmt.Errorf("Failed to Generate Session")
	}

	var accessSecretKey string

	accessSecretKey, err = obj.generateSecretKey()

	if err != nil {
		log.Println(err)
		return nil, nil, "", fmt.Errorf("Failed to Generate Session")
	}

	var sessionKey string = fmt.Sprintf("auth-token-%s-%d", role, userID)
	var sessionData auth.Session = auth.Session{
		AccessUUID:      tokenUUID.GetValue(),
		RefreshUUID:     refreshTokenUUID.GetValue(),
		IsPhoneVerified: isPhoneVerified,
	}

	err = obj.sessionRepository.SetSession(sessionKey, &sessionData, time.Hour*8760)

	if err != nil {
		return nil, nil, "", err
	}

	return token, refreshToken, accessSecretKey, nil

}

func (obj *authToken) ValidateToken(token *domain.JWT, isRefreshToken bool) (*auth.ValidateAuthTokenResponse, error) {

	var err error
	var keyExist bool

	var tokenPayload map[string]interface{} = token.GetPayload()
	var role string
	var userID float64
	var tokenUUID string

	role, keyExist = tokenPayload["role"].(string)

	if !keyExist {
		return nil, fmt.Errorf("Unauthorized")
	}

	userID, keyExist = tokenPayload["user_id"].(float64)

	if !keyExist {
		return nil, fmt.Errorf("Unauthorized")
	}

	tokenUUID, keyExist = tokenPayload["uuid"].(string)

	if !keyExist {
		return nil, fmt.Errorf("Unauthorized")
	}

	var sessionKey string = fmt.Sprintf("auth-token-%s-%d", role, int(userID))

	var sessionData *auth.Session

	sessionData, err = obj.sessionRepository.GetSession(sessionKey)

	if err != nil {
		return nil, fmt.Errorf("Unauthorized")
	}

	var cacheTokenUUID string

	if isRefreshToken {
		cacheTokenUUID = sessionData.RefreshUUID
	} else {
		cacheTokenUUID = sessionData.AccessUUID
	}

	if tokenUUID != cacheTokenUUID {
		return nil, fmt.Errorf("Unauthorized")
	}

	err = obj.sessionRepository.ExtendSessionExpiration(sessionKey, time.Hour*8760)

	if err != nil {
		return nil, fmt.Errorf("Unauthorized")
	}

	var response auth.ValidateAuthTokenResponse = auth.ValidateAuthTokenResponse{
		Role:            role,
		IsPhoneVerified: sessionData.IsPhoneVerified,
	}

	return &response, nil

}

func (obj *authToken) RegenerateAuthToken(refreshToken *domain.JWT) (*domain.JWT, string, error) {

	var err error
	var validateResponse *auth.ValidateAuthTokenResponse

	validateResponse, err = obj.ValidateToken(refreshToken, true)

	if err != nil {
		return nil, "", fmt.Errorf("Refresh Token Invalid")
	}

	var tokenPayload map[string]interface{} = refreshToken.GetPayload()
	var keyExist bool
	var role string
	var userID float64
	var tokenUUID string

	role, keyExist = tokenPayload["role"].(string)

	if !keyExist {
		return nil, "", fmt.Errorf("Refresh Token Invalid")
	}

	userID, keyExist = tokenPayload["user_id"].(float64)

	if !keyExist {
		return nil, "", fmt.Errorf("Refresh Token Invalid")
	}

	tokenUUID, keyExist = tokenPayload["uuid"].(string)

	if !keyExist {
		return nil, "", fmt.Errorf("Refresh Token Invalid")
	}

	var newAccessTokenUUID *domain.UUID = domain.NewUUID()

	var newAccessTokenPayload jwt.MapClaims = jwt.MapClaims{}
	newAccessTokenPayload["user_id"] = userID
	newAccessTokenPayload["role"] = role
	newAccessTokenPayload["uuid"] = newAccessTokenUUID.GetValue()
	newAccessTokenPayload["exp"] = time.Now().Add(time.Hour * 24).Unix()

	var newAccessToken *domain.JWT
	newAccessToken, err = domain.NewJWTFromPayload(newAccessTokenPayload, obj.secretKey)

	if err != nil {
		log.Println(err)
		return nil, "", fmt.Errorf("Failed to generate new session")
	}

	var accessSecretKey string

	accessSecretKey, err = obj.generateSecretKey()

	if err != nil {
		log.Println(err)
		return nil, "", fmt.Errorf("Failed to generate new session")
	}

	var sessionKey string = fmt.Sprintf("auth-token-%s-%d", role, int(userID))
	var sessionData auth.Session = auth.Session{
		AccessUUID:      newAccessTokenUUID.GetValue(),
		RefreshUUID:     tokenUUID,
		IsPhoneVerified: validateResponse.IsPhoneVerified,
	}

	err = obj.sessionRepository.SetSession(sessionKey, &sessionData, 0)

	if err != nil {
		return nil, "", err
	}

	return newAccessToken, accessSecretKey, nil

}

func (obj *authToken) RemoveAuthToken(token *domain.JWT) error {

	var err error
	var keyExist bool

	var tokenPayload map[string]interface{} = token.GetPayload()
	var role string
	var userID float64

	userID, keyExist = tokenPayload["user_id"].(float64)

	if !keyExist {
		return fmt.Errorf("Access Token Invalid")
	}

	role, keyExist = tokenPayload["role"].(string)

	if !keyExist {
		return fmt.Errorf("Access Token Invalid")
	}

	var sessionKey string = fmt.Sprintf("auth-token-%s-%d", role, int(userID))

	err = obj.sessionRepository.RemoveSession(sessionKey)

	if err != nil {
		return err
	}

	return nil

}
