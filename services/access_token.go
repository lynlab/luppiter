package services

import (
	"crypto/rand"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/hellodhlyn/luppiter/models"
	"github.com/hellodhlyn/luppiter/repositories"
)

type AccessTokenService interface {
	CreateAccessToken(identity *models.UserIdentity, app *models.Application) (*models.AccessToken, error)
	ActivateAccessToken(activationToken string) (*models.AccessToken, error)
}

type AccessTokenServiceImpl struct {
	repo repositories.AccessTokenRepository
}

func NewAccessTokenService(repo repositories.AccessTokenRepository) (AccessTokenService, error) {
	return &AccessTokenServiceImpl{repo}, nil
}

func (svc *AccessTokenServiceImpl) CreateAccessToken(identity *models.UserIdentity, app *models.Application) (*models.AccessToken, error) {
	token := &models.AccessToken{
		IdentityID:    identity.ID,
		Identity:      *identity,
		ApplicationID: app.ID,
		Application:   *app,
		AccessKey:     secureRandomString(20),
		SecretKey:     secureRandomString(20),
		ActivationKey: secureRandomString(20),
		Activated:     false,
	}

	svc.repo.Save(token)
	return token, nil
}

func (svc *AccessTokenServiceImpl) ActivateAccessToken(activationToken string) (*models.AccessToken, error) {
	token, _ := jwt.Parse(activationToken, nil)
	activationKey := token.Claims.(jwt.MapClaims)["activationKey"].(string)
	accessToken := svc.repo.FindByActivationKey(activationKey)
	if accessToken == nil {
		return nil, errors.New("access token not found")
	}

	token, err := jwt.Parse(activationToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(accessToken.Application.SecretKey), nil
	})
	if err != nil {
		return nil, errors.New("invalid token")
	}

	expireAt := time.Now().AddDate(0, 0, 7)
	accessToken.Activated = true
	accessToken.ExpireAt = &expireAt
	svc.repo.Save(accessToken)

	return accessToken, nil
}

func secureRandomString(l int) string {
	bytes := make([]byte, l)
	_, _ = rand.Read(bytes)
	return fmt.Sprintf("%x", bytes)
}
