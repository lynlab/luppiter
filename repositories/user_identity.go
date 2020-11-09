package repositories

import (
	"gorm.io/gorm"

	"github.com/hellodhlyn/luppiter/models"
)

type UserIdentityRepository interface {
	Save(identity *models.UserIdentity)
}

type UserIdentityRepositoryImpl struct {
	db *gorm.DB
}

func NewUserIdentityRepository(db *gorm.DB) (UserIdentityRepository, error) {
	return &UserIdentityRepositoryImpl{db}, nil
}

func (repo *UserIdentityRepositoryImpl) Save(identity *models.UserIdentity) {
	repo.db.Save(identity)
}
