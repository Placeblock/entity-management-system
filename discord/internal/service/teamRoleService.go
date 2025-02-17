package service

import (
	"context"

	teamrole "github.com/Placeblock/nostalgicraft-discord/internal/repository/teamRole"
	"github.com/Placeblock/nostalgicraft-discord/pkg/models"
)

type TeamRoleService struct {
	repo *teamrole.TeamRoleRepository
}

func NewTeamRoleService(repository *teamrole.TeamRoleRepository) *TeamRoleService {
	return &TeamRoleService{repo: repository}
}

func (service *TeamRoleService) GetRoleByTeamId(ctx context.Context, teamId uint) (string, error) {
	return (*service.repo).GetRoleByTeamId(ctx, teamId)
}

func (service *TeamRoleService) CreateTeamRole(ctx context.Context, teamId uint, roleId string) error {
	return (*service.repo).CreateTeamRole(ctx, models.TeamRole{TeamID: teamId, RoleID: roleId})
}

func (service *TeamRoleService) DeleteTeamRole(ctx context.Context, teamId uint) error {
	return (*service.repo).DeleteTeamRole(ctx, teamId)
}
