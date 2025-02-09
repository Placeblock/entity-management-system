package team

import (
	"context"

	"github.com/codelix/ems/pkg/models"
)

type TeamRepository interface {
	GetTeams(ctx context.Context, filter models.Team) (*[]models.Team, error)
	GetTeam(ctx context.Context, team *models.Team) error

	CreateTeam(ctx context.Context, team *models.Team, member *models.Member) error
	DeleteTeam(ctx context.Context, team *models.Team) error
	UpdateTeam(ctx context.Context, team models.Team) error
}
