package teamentity

import (
	"context"

	"github.com/codelix/ems/pkg/models"
)

type TeamEntityRepository interface {
	GetTeamEntities(ctx context.Context) (*[]models.TeamEntity, error)
	GetTeamEntitiesByTeamId(ctx context.Context, teamId uint) (*[]models.Entity, error)
	CreateTeamEntity(ctx context.Context, teamEntity *models.TeamEntity) error
	DeleteTeamEntity(ctx context.Context, entityId uint) error
	GetTeamEntityByEntityId(ctx context.Context, entityId uint) (*models.TeamEntity, error)
}
