package team

import (
	"context"

	"github.com/codelix/ems/pkg/models"
	"github.com/google/uuid"
)

type TeamRepository interface {
	GetTeams(ctx context.Context) (*[]models.Team, error)
	GetTeam(ctx context.Context, team *models.Team) error

	CreateTeam(ctx context.Context, team *models.Team) error
	DeleteTeam(ctx context.Context, id uuid.UUID) error
	UpdateTeam(ctx context.Context, team models.Team) error
}
