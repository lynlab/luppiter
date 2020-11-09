package services

import (
	"context"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/hellodhlyn/luppiter/models"
	"github.com/hellodhlyn/luppiter/repositories"
	"google.golang.org/api/idtoken"
)

const (
	providerGoogle = "google"
)

type UserAccountService interface {
	FindOrCreateByGoogleAccount(string) (*models.UserAccount, error)
}

type UserAccountServiceImpl struct {
	accountRepo  repositories.UserAccountRepository
	identityRepo repositories.UserIdentityRepository
	validator    *idtoken.Validator
	audience     string
}

func NewUserAccountService(accountRepo repositories.UserAccountRepository, identityRepo repositories.UserIdentityRepository) (UserAccountService, error) {
	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	googleCredPath := os.Getenv("GOOGLE_SECRET_ACCOUNT_PATH")

	validator, err := idtoken.NewValidator(context.Background(), idtoken.WithCredentialsFile(googleCredPath))
	if err != nil {
		return nil, err
	}

	return &UserAccountServiceImpl{accountRepo, identityRepo, validator, googleClientID}, nil
}

func (svc *UserAccountServiceImpl) FindOrCreateByGoogleAccount(idToken string) (*models.UserAccount, error) {
	payload, err := svc.validator.Validate(context.Background(), idToken, svc.audience)
	if err != nil {
		return nil, err
	}

	account := svc.accountRepo.FindByProviderId(providerGoogle, payload.Subject)
	if account == nil {
		// payload.Claims returns an empty map by a bug.
		// See: https://github.com/googleapis/google-api-go-client/pull/498
		token, _ := jwt.Parse(idToken, nil)
		claims := token.Claims.(jwt.MapClaims)

		identity := &models.UserIdentity{UUID: uuid.New().String(), Username: claims["name"].(string), Email: claims["email"].(string)}
		svc.identityRepo.Save(identity)

		account = &models.UserAccount{Provider: providerGoogle, ProviderID: payload.Subject, IdentityID: identity.ID, Identity: *identity}
		svc.accountRepo.Save(account)
	}

	return account, nil
}
