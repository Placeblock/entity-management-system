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

func (service *TeamEntityService) CreateInvite(ctx context.Context, invitedId uint, inviterId uint, teamId uint) error {
	teamEntityInvite := models.TeamEntityInvite{InvitedID: invitedId, InviterID: inviterId, TeamID: teamId}
	return (*service.teamEntityRepository).CreateTeamEntityInvite(ctx, teamEntityInvite)
}

func (service *TeamEntityService) ProcessInvite(ctx context.Context, invitedId uint, teamId uint, accept bool) error {
	teamEntityInvite, err := (*service.teamEntityRepository).GetTeamEntityInvite(ctx, invitedId, teamId)
	if err != nil {
		return err
	}
	if accept {
		return (*service.teamEntityRepository).AcceptTeamEntityInvite(ctx, *teamEntityInvite)
	} else {
		return (*service.teamEntityRepository).DeclineTeamEntityInvite(ctx, *teamEntityInvite)
	}
}

func (service *TeamEntityService) GetTeamEntityInvites(ctx context.Context) (*[]models.TeamEntityInvite, error) {
	return (*service.teamEntityRepository).GetTeamEntityInvites(ctx)
}

func (service *TeamEntityService) GetTeamEntityInvitesByInvitedId(ctx context.Context, invitedId uint) (*[]models.TeamEntityInvite, error) {
	return (*service.teamEntityRepository).GetTeamEntityInvitesByInvitedId(ctx, invitedId)
}
