package teamentity

import (
	"context"

	"github.com/codelix/ems/internal/realtime"
	teamentity "github.com/codelix/ems/internal/repository/teamEntity"
	"github.com/codelix/ems/pkg/models"
)

type TeamEntityService struct {
	teamEntityRepository *teamentity.TeamEntityRepository
	publisher            *realtime.Publisher
}

func NewTeamEntityService(repository teamentity.TeamEntityRepository, publisher *realtime.Publisher) *TeamEntityService {
	return &TeamEntityService{&repository, publisher}
}

func (service *TeamEntityService) SetTeam(ctx context.Context, entityId uint, teamId uint) error {
	teamEntity := models.TeamEntity{EntityID: entityId, TeamID: teamId}
	return (*service.teamEntityRepository).CreateTeamEntity(ctx, &teamEntity)
}

func (service *TeamEntityService) LeaveTeam(ctx context.Context, entityId uint) error {
	return (*service.teamEntityRepository).DeleteTeamEntity(ctx, entityId)
}

func (service *TeamEntityService) GetTeamEntities(ctx context.Context) (*[]models.TeamEntity, error) {
	return (*service.teamEntityRepository).GetTeamEntities(ctx)
}

func (service *TeamEntityService) GetTeamEntitiesByTeamId(ctx context.Context, teamId uint) (*[]models.Entity, error) {
	return (*service.teamEntityRepository).GetTeamEntitiesByTeamId(ctx, teamId)
}

func (service *TeamEntityService) GetTeamEntity(ctx context.Context, entityId uint) (*models.TeamEntity, error) {
	return (*service.teamEntityRepository).GetTeamEntityByEntityId(ctx, entityId)
}
