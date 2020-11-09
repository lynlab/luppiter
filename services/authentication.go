package services

import (
	"errors"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/hellodhlyn/luppiter/models"
	"github.com/hellodhlyn/luppiter/repositories"
)

type AuthenticationService interface {
	Authenticate(*http.Request) (*models.UserIdentity, error)
}

type AuthenticationServiceImpl struct {
	tokenRepo repositories.AccessTokenRepository
}

func NewAuthenticationService(tokenRepo repositories.AccessTokenRepository) (AuthenticationService, error) {
	return &AuthenticationServiceImpl{tokenRepo}, nil
}

func (svc *AuthenticationServiceImpl) Authenticate(r *http.Request) (*models.UserIdentity, error) {
	authorization := r.Header.Get("Authorization")
	splits := strings.Split(authorization, " ")
	if len(splits) != 2 {
		return nil, errors.New("invalid authorization")
	}
	jwtString := splits[len(splits)-1]

	token, _ := jwt.Parse(jwtString, nil)
	accessKey := token.Claims.(jwt.MapClaims)["accessKey"].(string)
	accessToken := svc.tokenRepo.FindByAccessKey(accessKey)
	if accessToken == nil {
		return nil, errors.New("invalid access key")
	}

	token, err := jwt.Parse(jwtString, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(accessToken.SecretKey), nil
	})
	if err != nil {
		return nil, errors.New("invalid signature")
	}
	if accessToken.HasExpired() {
		return nil, errors.New("access token expired")
	}

	return &accessToken.Identity, nil
}
