package teamrole

import (
	"context"

	"github.com/Placeblock/nostalgicraft-discord/pkg/models"
)

type TeamRoleRepository interface {
	GetRoleByTeamId(ctx context.Context, teamId uint) (string, error)
	CreateTeamRole(ctx context.Context, teamRole models.TeamRole) error
	DeleteTeamRole(ctx context.Context, teamId uint) error
}
