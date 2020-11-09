package repository

import (
	"context"
	"net/http"
	"os"

	"google.golang.org/api/idtoken"
	"gorm.io/gorm"

	"github.com/hellodhlyn/luppiter/model"
)

type UserAccountRepository interface {
	FindByProviderId(string, string) *model.UserAccount
	Save(*model.UserAccount)
}

type UserAccountRepositoryImpl struct {
	idTknClient *http.Client
	db          *gorm.DB
}

func NewUserAccountRepository(db *gorm.DB) (UserAccountRepository, error) {
	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	googleCredPath := os.Getenv("GOOGLE_SECRET_ACCOUNT_PATH")
	idTknClient, err := idtoken.NewClient(context.Background(), googleClientID, idtoken.WithCredentialsFile(googleCredPath))
	if err != nil {
		return nil, err
	}

	return &UserAccountRepositoryImpl{idTknClient, db}, nil
}

func (repo *UserAccountRepositoryImpl) FindByProviderId(provider, providerId string) *model.UserAccount {
	var account model.UserAccount
	repo.db.Where(&model.UserAccount{Provider: provider, ProviderID: providerId}).Preload("Identity").First(&account)
	if account.ID == 0 {
		return nil
	}
	return &account
}

func (repo *UserAccountRepositoryImpl) Save(account *model.UserAccount) {
	repo.db.Save(account)
}
