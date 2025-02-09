package team

import (
	"context"

	"github.com/codelix/ems/pkg/models"
)

type TeamRepository interface {
	GetTeams(ctx context.Context) (*[]models.Team, error)
	GetTeam(ctx context.Context, team *models.Team) error

	CreateTeam(ctx context.Context, team *models.Team) error
	DeleteTeam(ctx context.Context, id uint) error
	UpdateTeam(ctx context.Context, team models.Team) error
}
