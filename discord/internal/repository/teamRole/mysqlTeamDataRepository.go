package teamrole

import (
	"context"
	"errors"

	cerr "github.com/Placeblock/nostalgicraft-discord/pkg/errors"
	"github.com/Placeblock/nostalgicraft-discord/pkg/models"
	"gorm.io/gorm"
)

type MysqlTeamDataRepository struct {
	db *gorm.DB
}

func NewMysqlTeamDataRepository(db *gorm.DB) *MysqlTeamDataRepository {
	return &MysqlTeamDataRepository{db}
}

func (repo *MysqlTeamDataRepository) GetTeamData(ctx context.Context, teamData *models.TeamData) error {
	if err := repo.db.WithContext(ctx).First(&teamData).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &cerr.ErrNotFound{}
		}
		return err
	}
	return nil
}

func (repo *MysqlTeamDataRepository) CreateTeamData(ctx context.Context, teamRole models.TeamData) error {
	return repo.db.WithContext(ctx).Create(&teamRole).Error
}

func (repo *MysqlTeamDataRepository) DeleteTeamData(ctx context.Context, teamId uint) error {
	teamRole := models.TeamData{TeamID: teamId}
	return repo.db.WithContext(ctx).Delete(&teamRole).Error
}
