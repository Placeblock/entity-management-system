package team

import (
	"context"

	"github.com/codelix/ems/internal/realtime"
	"github.com/codelix/ems/internal/repository/team"
	"github.com/codelix/ems/pkg/models"
)

type TeamService struct {
	teamRepository *team.TeamRepository
	publisher      *realtime.Publisher
}

func NewMysqlTeamRepository(repo team.TeamRepository, publisher *realtime.Publisher) *TeamService {
	return &TeamService{&repo, publisher}
}

func (service *TeamService) GetTeams(ctx context.Context) (*[]models.Team, error) {
	return (*service.teamRepository).GetTeams(ctx)
}

func (service *TeamService) GetTeam(ctx context.Context, teamId uint) (*models.Team, error) {
	team := models.Team{ID: teamId}
	err := (*service.teamRepository).GetTeam(ctx, &team)
	if err != nil {
		return nil, err
	}
	return &team, nil
}

func (service *TeamService) CreateTeam(ctx context.Context, team *models.Team) error {
	err := (*service.teamRepository).CreateTeam(ctx, team)
	if err != nil {
		return err
	}
	service.publisher.Channel <- realtime.Action{Type: "team.create", Data: team}
	return nil
}

func (service *TeamService) RenameTeam(ctx context.Context, id uint, newName string) error {
	team := models.Team{ID: id, Name: newName}
	err := (*service.teamRepository).UpdateTeam(ctx, team)
	if err != nil {
		return err
	}
	err = (*service.teamRepository).GetTeam(ctx, &team)
	if err != nil {
		return err
	}
	service.publisher.Channel <- realtime.Action{Type: "team.rename", Data: team}
	return nil
}

func (service *TeamService) RecolorTeam(ctx context.Context, id uint, newHue models.Hue) error {
	team := models.Team{ID: id, Hue: &newHue}
	err := (*service.teamRepository).UpdateTeam(ctx, team)
	if err != nil {
		return err
	}
	err = (*service.teamRepository).GetTeam(ctx, &team)
	if err != nil {
		return err
	}
	service.publisher.Channel <- realtime.Action{Type: "team.recolor", Data: team}
	return nil
}

func (service *TeamService) SetOwner(ctx context.Context, id uint, newOwner uint) error {
	team := models.Team{ID: id, OwnerID: newOwner}
	err := (*service.teamRepository).UpdateTeam(ctx, team)
	if err != nil {
		return err
	}
	err = (*service.teamRepository).GetTeam(ctx, &team)
	if err != nil {
		return err
	}
	service.publisher.Channel <- realtime.Action{Type: "team.owner", Data: team}
	return nil
}
