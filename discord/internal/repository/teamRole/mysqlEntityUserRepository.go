package teamrole

import (
	"context"
	"errors"

	cerr "github.com/Placeblock/nostalgicraft-discord/pkg/errors"
	"github.com/Placeblock/nostalgicraft-discord/pkg/models"
	"gorm.io/gorm"
)

type MysqlTeamRoleRepository struct {
	db *gorm.DB
}

func NewMysqlTeamRoleRepository(db *gorm.DB) *MysqlTeamRoleRepository {
	return &MysqlTeamRoleRepository{db}
}

func (repo *MysqlTeamRoleRepository) GetRoleByTeamId(ctx context.Context, teamId uint) (string, error) {
	teamRole := models.TeamRole{TeamID: teamId}
	if err := repo.db.WithContext(ctx).First(&teamRole).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", &cerr.ErrNotFound{}
		}
		return "", err
	}
	return teamRole.RoleID, nil
}

func (repo *MysqlTeamRoleRepository) CreateTeamRole(ctx context.Context, teamRole models.TeamRole) error {
	return repo.db.WithContext(ctx).Create(&teamRole).Error
}

func (repo *MysqlTeamRoleRepository) DeleteTeamRole(ctx context.Context, teamId uint) error {
	teamRole := models.TeamRole{TeamID: teamId}
	return repo.db.WithContext(ctx).Delete(&teamRole).Error
}
