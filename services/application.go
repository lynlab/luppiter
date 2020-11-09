package services

import (
	"github.com/hellodhlyn/luppiter/models"
	"github.com/hellodhlyn/luppiter/repositories"
)

type ApplicationService interface {
	FindByUUID(uuid string) *models.Application
}

type ApplicationServiceImpl struct {
	repo repositories.ApplicationRepository
}

func NewApplicationService(repo repositories.ApplicationRepository) (ApplicationService, error) {
	return &ApplicationServiceImpl{repo}, nil
}

func (svc *ApplicationServiceImpl) FindByUUID(uuid string) *models.Application {
	return svc.repo.FindByUUID(uuid)
}
