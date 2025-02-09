package entity

import (
	"context"

	"github.com/codelix/ems/pkg/models"
)

type EntityRepository interface {
	GetEntity(ctx context.Context, entity *models.Entity) error
	CreateEntity(ctx context.Context, entity *models.Entity) error
	UpdateEntity(ctx context.Context, entity models.Entity) error
	GetEntities(ctx context.Context, filter models.Entity) (*[]models.Entity, error)
}
