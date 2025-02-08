package team

import (
	"context"

	"github.com/codelix/ems/internal/repository/team"
	"github.com/codelix/ems/pkg/models"
)

type TeamService struct {
	teamRepository *team.TeamRepository
}

func NewMysqlTeamRepository(repo *team.TeamRepository) *TeamService {
	return &TeamService{repo}
}

func (service *TeamService) GetTeams(ctx context.Context) (*[]models.Team, error) {
	return (*service.teamRepository).GetTeams(ctx)
}

func (service *TeamService) GetTeam(ctx context.Context, teamId uint) (*models.Team, error) {
	var team models.Team
	err := (*service.teamRepository).GetTeam(ctx, &team)
	if err != nil {
		return nil, err
	}
	return &team, nil
}

func (service *TeamService) CreateTeam(ctx context.Context, team *models.Team) error {
	return (*service.teamRepository).CreateTeam(ctx, team)
}

func (service *TeamService) RenameTeam(ctx context.Context, id uint, newName string) error {
	team := models.Team{ID: id, Name: newName}
	return (*service.teamRepository).UpdateTeam(ctx, team)
}

func (service *TeamService) RecolorTeam(ctx context.Context, id uint, newHue models.Hue) error {
	team := models.Team{ID: id, Hue: newHue}
	return (*service.teamRepository).UpdateTeam(ctx, team)
}

func (service *TeamService) SetOwner(ctx context.Context, id uint, newOwner uint) error {
	team := models.Team{ID: id, OwnerID: newOwner}
	return (*service.teamRepository).UpdateTeam(ctx, team)
}
