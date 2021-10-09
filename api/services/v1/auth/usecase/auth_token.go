package usecase

import (
	"fmt"
	"log"
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

func (obj *authToken) GenerateAuthToken(userID string, role string, isPhoneVerified bool) (*domain.JWT, *domain.JWT, error) {

	var err error

	var token, refreshToken *domain.JWT
	var tokenPayload, refreshTokenPayload jwt.MapClaims
	var tokenUUID *domain.UUID = domain.NewUUID()
	var refreshTokenUUID *domain.UUID = domain.NewUUID()

	tokenPayload = jwt.MapClaims{}
	tokenPayload["user_id"] = userID
	tokenPayload["role"] = role
	tokenPayload["uuid"] = tokenUUID.GetValue()
	tokenPayload["exp"] = time.Now().Add(time.Second * 10).Unix()

	token, err = domain.NewJWTFromPayload(tokenPayload, obj.secretKey)

	if err != nil {
		log.Println(err)
		return nil, nil, fmt.Errorf("Failed to Generate Session")
	}

	refreshTokenPayload = jwt.MapClaims{}
	refreshTokenPayload["user_id"] = userID
	refreshTokenPayload["role"] = role
	refreshTokenPayload["uuid"] = refreshTokenUUID.GetValue()

	refreshToken, err = domain.NewJWTFromPayload(refreshTokenPayload, obj.secretKey)

	if err != nil {
		log.Println(err)
		return nil, nil, fmt.Errorf("Failed to Generate Session")
	}

	token, err = domain.NewJWTFromPayload(tokenPayload, obj.secretKey)

	if err != nil {
		log.Println(err)
		return nil, nil, fmt.Errorf("Failed to Generate Session")
	}

	var sessionKey string = fmt.Sprintf("auth-token-%s-%s", role, userID)
	var sessionData auth.Session = auth.Session{
		AccessUUID:      tokenUUID.GetValue(),
		RefreshUUID:     refreshTokenUUID.GetValue(),
		IsPhoneVerified: isPhoneVerified,
	}

	err = obj.sessionRepository.SetSession(sessionKey, &sessionData, time.Hour*8760)

	if err != nil {
		return nil, nil, err
	}

	return token, refreshToken, nil

}

func (obj *authToken) ValidateToken(token *domain.JWT, isRefreshToken bool) (*auth.ValidateAuthTokenResponse, error) {

	var err error
	var keyExist bool

	var tokenPayload map[string]interface{} = token.GetPayload()
	var role string
	var userID string
	var tokenUUID string

	role, keyExist = tokenPayload["role"].(string)

	if !keyExist {
		return nil, fmt.Errorf("Unauthorized")
	}

	userID, keyExist = tokenPayload["user_id"].(string)

	if !keyExist {
		return nil, fmt.Errorf("Unauthorized")
	}

	tokenUUID, keyExist = tokenPayload["uuid"].(string)

	if !keyExist {
		return nil, fmt.Errorf("Unauthorized")
	}

	var sessionKey string = fmt.Sprintf("auth-token-%s-%s", role, userID)

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

	var userUUID *domain.UUID

	userUUID, err = domain.NewUUIDFromString(userID)

	if err != nil {
		return nil, fmt.Errorf("Unauthorized")
	}

	var response auth.ValidateAuthTokenResponse = auth.ValidateAuthTokenResponse{
		UserID:          userUUID,
		Role:            role,
		IsPhoneVerified: sessionData.IsPhoneVerified,
	}

	return &response, nil

}

func (obj *authToken) RegenerateAuthToken(refreshToken *domain.JWT) (*domain.JWT, error) {

	var err error
	var validateResponse *auth.ValidateAuthTokenResponse

	validateResponse, err = obj.ValidateToken(refreshToken, true)

	if err != nil {
		return nil, fmt.Errorf("Refresh Token Invalid")
	}

	var tokenPayload map[string]interface{} = refreshToken.GetPayload()
	var keyExist bool
	var role string
	var userID string
	var tokenUUID string

	role, keyExist = tokenPayload["role"].(string)

	if !keyExist {
		return nil, fmt.Errorf("Refresh Token Invalid")
	}

	userID, keyExist = tokenPayload["user_id"].(string)

	if !keyExist {
		return nil, fmt.Errorf("Refresh Token Invalid")
	}

	tokenUUID, keyExist = tokenPayload["uuid"].(string)

	if !keyExist {
		return nil, fmt.Errorf("Refresh Token Invalid")
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
		return nil, fmt.Errorf("Failed to generate new session")
	}

	var sessionKey string = fmt.Sprintf("auth-token-%s-%s", role, userID)
	var sessionData auth.Session = auth.Session{
		AccessUUID:      newAccessTokenUUID.GetValue(),
		RefreshUUID:     tokenUUID,
		IsPhoneVerified: validateResponse.IsPhoneVerified,
	}

	err = obj.sessionRepository.SetSession(sessionKey, &sessionData, 0)

	if err != nil {
		return nil, err
	}

	return newAccessToken, nil

}

func (obj *authToken) RemoveAuthToken(token *domain.JWT) error {

	var err error
	var keyExist bool

	var tokenPayload map[string]interface{} = token.GetPayload()
	var role string
	var userID string

	userID, keyExist = tokenPayload["user_id"].(string)

	if !keyExist {
		return fmt.Errorf("Access Token Invalid")
	}

	role, keyExist = tokenPayload["role"].(string)

	if !keyExist {
		return fmt.Errorf("Access Token Invalid")
	}

	var sessionKey string = fmt.Sprintf("auth-token-%s-%s", role, userID)

	err = obj.sessionRepository.RemoveSession(sessionKey)

	if err != nil {
		return err
	}

	return nil

}
