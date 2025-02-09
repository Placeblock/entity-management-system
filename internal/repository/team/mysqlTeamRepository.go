package team

import (
	"context"
	"fmt"

	"github.com/codelix/ems/pkg/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MysqlTeamRepository struct {
	db *gorm.DB
}

func NewMysqlTeamRepository(db *gorm.DB) *MysqlTeamRepository {
	return &MysqlTeamRepository{db}
}

func (repo *MysqlTeamRepository) GetTeams(ctx context.Context, filter models.Team) (*[]models.Team, error) {
	var teams []models.Team
	if err := repo.db.WithContext(ctx).Where(filter).Preload(clause.Associations).Find(&teams).Error; err != nil {
		return nil, fmt.Errorf("getTeams %+v: %v", filter, err.Error())
	}
	return &teams, nil
}

func (repo *MysqlTeamRepository) GetTeam(ctx context.Context, team *models.Team) error {
	return repo.db.WithContext(ctx).Preload(clause.Associations).First(&team, &team.ID).Error
}

func (repo *MysqlTeamRepository) CreateTeam(ctx context.Context, team *models.Team) error {
	if err := repo.db.WithContext(ctx).Create(team).Preload(clause.Associations).Find(team).Error; err != nil {
		return fmt.Errorf("createTeam %+v: %v", team, err)
	}
	return nil
}

func (repo *MysqlTeamRepository) DeleteTeam(ctx context.Context, team *models.Team) error {
	if err := repo.db.WithContext(ctx).Delete(team).Error; err != nil {
		return fmt.Errorf("deleteTeam %+v: %v", team, err)
	}
	return nil
}

func (repo *MysqlTeamRepository) UpdateTeam(ctx context.Context, team models.Team) error {
	if err := repo.db.WithContext(ctx).Updates(&team).Error; err != nil {
		return fmt.Errorf("updateTeam %d: %v", team.ID, err)
	}
	return nil
}
