package teamrole

import (
	"context"

	"github.com/Placeblock/nostalgicraft-discord/pkg/models"
)

type TeamRoleRepository interface {
	GetTeamData(ctx context.Context, teamData *models.TeamData) error
	CreateTeamData(ctx context.Context, teamData models.TeamData) error
	DeleteTeamData(ctx context.Context, teamId uint) error
}
