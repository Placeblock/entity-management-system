package team

import (
	"context"

	"github.com/Placeblock/nostalgicraft-ems/pkg/models"
)

type TeamRepository interface {
	GetTeams(ctx context.Context, filter models.Team) (*[]models.Team, error)
	GetTeam(ctx context.Context, team *models.Team) error
	GetTeamByEntityID(ctx context.Context, entityId uint) (*models.Team, error)
	CreateTeam(ctx context.Context, team *models.Team, member *models.Member) error
	DeleteTeam(ctx context.Context, team *models.Team) error
	UpdateTeam(ctx context.Context, team models.Team) error
}
