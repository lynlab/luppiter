package repositories

import (
	"context"
	"net/http"
	"os"

	"google.golang.org/api/idtoken"
	"gorm.io/gorm"

	"github.com/hellodhlyn/luppiter/models"
)

type UserAccountRepository interface {
	FindByProviderId(string, string) *models.UserAccount
	Save(*models.UserAccount)
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

func (repo *UserAccountRepositoryImpl) FindByProviderId(provider, providerId string) *models.UserAccount {
	var account models.UserAccount
	repo.db.Where(&models.UserAccount{Provider: provider, ProviderID: providerId}).Preload("Identity").First(&account)
	if account.ID == 0 {
		return nil
	}
	return &account
}

func (repo *UserAccountRepositoryImpl) Save(account *models.UserAccount) {
	repo.db.Save(account)
}
