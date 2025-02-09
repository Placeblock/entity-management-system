package entity

import (
	"context"

	"github.com/codelix/ems/internal/realtime"
	"github.com/codelix/ems/internal/repository/entity"
	"github.com/codelix/ems/pkg/models"
)

type EntityService struct {
	entityRepository *entity.EntityRepository
	publisher        *realtime.Publisher
}

func NewEntityService(repository entity.EntityRepository, publisher *realtime.Publisher) *EntityService {
	return &EntityService{&repository, publisher}
}

func (service *EntityService) CreateEntity(ctx context.Context, entity *models.Entity) error {
	return (*service.entityRepository).CreateEntity(ctx, entity)
}

func (service *EntityService) RenameEntity(ctx context.Context, id uint, newName string) error {
	entity, err := (*service.entityRepository).GetEntity(ctx, id)
	if err != nil {
		return err
	}
	entity.Name = newName
	(*service.entityRepository).UpdateEntity(ctx, *entity)
	service.publisher.Channel <- realtime.Action{Type: "entity.rename", Data: entity}
	return nil
}

func (service *EntityService) GetEntities(ctx context.Context) (*[]models.Entity, error) {
	return (*service.entityRepository).GetEntities(ctx, models.Entity{})
}

func (service *EntityService) GetEntity(ctx context.Context, id uint) (*models.Entity, error) {
	entity := models.Entity{ID: id}
	err := (*service.entityRepository).GetEntity(ctx, &entity)
	if err != nil {
		return nil, err
	}
	return &entity, nil
}
