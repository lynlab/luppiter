package repositories

import (
	"gorm.io/gorm"

	"github.com/hellodhlyn/luppiter/models"
)

type ApplicationRepository interface {
	FindByUUID(uuid string) *models.Application
}

type ApplicationRepositoryImpl struct {
	db *gorm.DB
}

func NewApplicationRepository(db *gorm.DB) (ApplicationRepository, error) {
	return &ApplicationRepositoryImpl{db}, nil
}

func (repo *ApplicationRepositoryImpl) FindByUUID(uuid string) *models.Application {
	var application models.Application
	repo.db.Where(&models.Application{UUID: uuid}).First(&application)
	if application.ID == 0 {
		return nil
	}
	return &application
}
