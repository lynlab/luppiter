package repository

import (
	"gorm.io/gorm"

	"github.com/hellodhlyn/luppiter/model"
)

type UserIdentityRepository interface {
	Save(identity *model.UserIdentity)
}

type UserIdentityRepositoryImpl struct {
	db *gorm.DB
}

func NewUserIdentityRepository(db *gorm.DB) (UserIdentityRepository, error) {
	return &UserIdentityRepositoryImpl{db}, nil
}

func (repo *UserIdentityRepositoryImpl) Save(identity *model.UserIdentity) {
	repo.db.Save(identity)
}
