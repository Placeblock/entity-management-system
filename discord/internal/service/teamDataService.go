package service

import (
	"context"
	"errors"

	teamrole "github.com/Placeblock/nostalgicraft-discord/internal/repository/teamRole"
	cerr "github.com/Placeblock/nostalgicraft-discord/pkg/errors"
	"github.com/Placeblock/nostalgicraft-discord/pkg/models"
)

type TeamDataService struct {
	repo *teamrole.TeamRoleRepository
}

func NewTeamDataService(repository teamrole.TeamRoleRepository) *TeamDataService {
	return &TeamDataService{repo: &repository}
}

func (service *TeamDataService) GetTeamDataByTeamId(ctx context.Context, teamId uint) (*models.TeamData, error) {
	teamData := models.TeamData{TeamID: teamId}
	err := (*service.repo).GetTeamData(ctx, &teamData)
	if err != nil {
		if (errors.Is(err, cerr.ErrNotFound{})) {
			return nil, nil
		}
		return nil, err
	}
	return &teamData, nil
}

func (service *TeamDataService) GetTeamDataByChannelId(ctx context.Context, channelId string) (*models.TeamData, error) {
	teamData := models.TeamData{ChannelID: channelId}
	err := (*service.repo).GetTeamData(ctx, &teamData)
	if err != nil {
		if (errors.Is(err, cerr.ErrNotFound{})) {
			return nil, nil
		}
		return nil, err
	}
	return &teamData, nil
}

func (service *TeamDataService) CreateTeamData(ctx context.Context, teamId uint, roleId string, channelId string) error {
	return (*service.repo).CreateTeamData(ctx, models.TeamData{TeamID: teamId, RoleID: roleId, ChannelID: channelId})
}

func (service *TeamDataService) DeleteTeamData(ctx context.Context, teamId uint) error {
	return (*service.repo).DeleteTeamData(ctx, teamId)
}
