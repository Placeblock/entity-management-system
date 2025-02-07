package team

import (
	"context"
	"fmt"

	"github.com/codelix/ems/pkg/models"
	"gorm.io/gorm"
)

type MysqlTeamRepository struct {
	db *gorm.DB
}

func NewMysqlTeamRepository(db *gorm.DB) *MysqlTeamRepository {
	return &MysqlTeamRepository{db}
}

func (repo *MysqlTeamRepository) GetTeams(ctx context.Context) (*[]models.Team, error) {
	var teams []models.Team
	if err := repo.db.WithContext(ctx).Preload("Owner").Find(&teams).Error; err != nil {
		return nil, fmt.Errorf("getTeams: %v", err.Error())
	}
	return &teams, nil
}

func (repo *MysqlTeamRepository) GetTeam(ctx context.Context, team *models.Team) error {
	return repo.db.WithContext(ctx).Preload("Owner").First(&team, &team.ID).Error
}

func (repo *MysqlTeamRepository) CreateTeam(ctx context.Context, team *models.Team) error {
	if err := repo.db.WithContext(ctx).Create(team).Error; err != nil {
		return fmt.Errorf("createTeam %s: %v", team.Name, err)
	}
	return nil
}

func (repo *MysqlTeamRepository) DeleteTeam(ctx context.Context, id uint) error {
	if err := repo.db.WithContext(ctx).Delete(models.Team{ID: id}).Error; err != nil {
		return fmt.Errorf("deleteTeam %d: %v", id, err)
	}
	return nil
}

func (repo *MysqlTeamRepository) UpdateTeam(ctx context.Context, team models.Team) error {
	if err := repo.db.WithContext(ctx).Save(team).Error; err != nil {
		return fmt.Errorf("updateTeam %d: %v", team.ID, err)
	}
	return nil
}
