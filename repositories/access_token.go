package repositories

import (
	"gorm.io/gorm"

	"github.com/hellodhlyn/luppiter/models"
)

type AccessTokenRepository interface {
	FindByAccessKey(string) *models.AccessToken
	FindByActivationKey(string) *models.AccessToken
	Save(*models.AccessToken)
}

type AccessTokenRepositoryImpl struct {
	db *gorm.DB
}

func NewAccessTokenRepository(db *gorm.DB) (AccessTokenRepository, error) {
	return &AccessTokenRepositoryImpl{db}, nil
}

func (repo *AccessTokenRepositoryImpl) FindByAccessKey(accessKey string) *models.AccessToken {
	var token models.AccessToken
	repo.db.Where(&models.AccessToken{AccessKey: accessKey}).Preload("Identity").Preload("Application").First(&token)
	if token.ID == 0 {
		return nil
	}
	return &token
}

func (repo *AccessTokenRepositoryImpl) FindByActivationKey(activationKey string) *models.AccessToken {
	var token models.AccessToken
	repo.db.Where(&models.AccessToken{ActivationKey: activationKey}).Preload("Identity").Preload("Application").First(&token)
	if token.ID == 0 {
		return nil
	}
	return &token
}

func (repo *AccessTokenRepositoryImpl) Save(token *models.AccessToken) {
	repo.db.Save(token)
}
