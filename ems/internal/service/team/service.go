package team

import (
	"context"

	"github.com/Placeblock/nostalgicraft-ems/internal/realtime"
	"github.com/Placeblock/nostalgicraft-ems/internal/repository/team"
	"github.com/Placeblock/nostalgicraft-ems/pkg/models"
	rtm "github.com/Placeblock/nostalgicraft-ems/pkg/realtime"
)

type TeamService struct {
	teamRepository *team.TeamRepository
	publisher      *realtime.Publisher
}

func NewMysqlTeamRepository(repo team.TeamRepository, publisher *realtime.Publisher) *TeamService {
	return &TeamService{&repo, publisher}
}

func (service *TeamService) GetTeams(ctx context.Context) (*[]models.Team, error) {
	return (*service.teamRepository).GetTeams(ctx, models.Team{})
}

func (service *TeamService) GetTeam(ctx context.Context, teamId uint) (*models.Team, error) {
	team := models.Team{ID: teamId}
	err := (*service.teamRepository).GetTeam(ctx, &team)
	if err != nil {
		return nil, err
	}
	return &team, nil
}

func (service *TeamService) GetTeamByEntityID(ctx context.Context, entityId uint) (*models.Team, error) {
	return (*service.teamRepository).GetTeamByEntityID(ctx, entityId)
}

func (service *TeamService) CreateTeam(ctx context.Context, team *models.Team, entityId uint) (*models.Member, error) {
	member := models.Member{EntityID: entityId}
	err := (*service.teamRepository).CreateTeam(ctx, team, &member)
	if err != nil {
		return nil, err
	}
	service.publisher.Channel <- rtm.Action{Type: "team.create", Data: rtm.CreateTeamData{Team: *team, Member: member}}
	return &member, nil
}

func (service *TeamService) RenameTeam(ctx context.Context, id uint, newName string) error {
	team := models.Team{ID: id}
	(*service.teamRepository).GetTeam(ctx, &team)
	team.Name = newName
	err := (*service.teamRepository).UpdateTeam(ctx, team)
	if err != nil {
		return err
	}
	service.publisher.Channel <- rtm.Action{Type: "team.rename", Data: team}
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
	service.publisher.Channel <- rtm.Action{Type: "team.recolor", Data: team}
	return nil
}
