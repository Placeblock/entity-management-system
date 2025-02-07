package entity

import (
	"context"

	"github.com/codelix/ems/internal/repository/entity"
	"github.com/codelix/ems/pkg/models"
)

type EntityService struct {
	entityRepository entity.EntityRepository
}

func NewEntityService(repository entity.EntityRepository) *EntityService {
	return &EntityService{repository}
}

func (service *EntityService) CreateEntity(ctx context.Context, entity *models.Entity) error {
	return service.entityRepository.CreateEntity(ctx, entity)
}

func (service *EntityService) RenameEntity(ctx context.Context, id int64, newName string) error {
	entity, err := service.entityRepository.GetEntity(ctx, id)
	if err != nil {
		return err
	}
	entity.Name = newName
	service.entityRepository.UpdateEntity(ctx, *entity)
	return nil
}

func (service *EntityService) GetEntities(ctx context.Context) (*[]models.Entity, error) {
	return service.entityRepository.GetEntities(ctx)
}

func (service *EntityService) GetEntity(ctx context.Context, id int64) (*models.Entity, error) {
	return service.entityRepository.GetEntity(ctx, id)
}
