package repository

import (
	"gorm.io/gorm"

	"github.com/hellodhlyn/luppiter/model"
)

type AccessTokenRepository interface {
	FindByAccessKey(string) *model.AccessToken
	FindByActivationKey(string) *model.AccessToken
	Save(*model.AccessToken)
}

type AccessTokenRepositoryImpl struct {
	db *gorm.DB
}

func NewAccessTokenRepository(db *gorm.DB) (AccessTokenRepository, error) {
	return &AccessTokenRepositoryImpl{db}, nil
}

func (repo *AccessTokenRepositoryImpl) FindByAccessKey(accessKey string) *model.AccessToken {
	var token model.AccessToken
	repo.db.Where(&model.AccessToken{AccessKey: accessKey}).Preload("Identity").Preload("Application").First(&token)
	if token.ID == 0 {
		return nil
	}
	return &token
}

func (repo *AccessTokenRepositoryImpl) FindByActivationKey(activationKey string) *model.AccessToken {
	var token model.AccessToken
	repo.db.Where(&model.AccessToken{ActivationKey: activationKey}).Preload("Identity").Preload("Application").First(&token)
	if token.ID == 0 {
		return nil
	}
	return &token
}

func (repo *AccessTokenRepositoryImpl) Save(token *model.AccessToken) {
	repo.db.Save(token)
}
