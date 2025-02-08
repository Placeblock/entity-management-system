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

	GetTeamEntityInvites(ctx context.Context) (*[]models.TeamEntityInvite, error)
	GetTeamEntityInvite(ctx context.Context, invitedId uint, teamId uint) (*models.TeamEntityInvite, error)
	GetTeamEntityInvitesByInvitedId(ctx context.Context, invitedId uint) (*[]models.TeamEntityInvite, error)
	CreateTeamEntityInvite(ctx context.Context, invite models.TeamEntityInvite) error
	DeclineTeamEntityInvite(ctx context.Context, teamEntityInvite models.TeamEntityInvite) error
	AcceptTeamEntityInvite(ctx context.Context, teamEntityInvite models.TeamEntityInvite) error
}
