package team

import (
	"context"
	"errors"
	"fmt"

	"github.com/Placeblock/nostalgicraft-ems/pkg/models"
	"gorm.io/gorm"
)

type MysqlTeamRepository struct {
	db *gorm.DB
}

func NewMysqlTeamRepository(db *gorm.DB) *MysqlTeamRepository {
	return &MysqlTeamRepository{db}
}

func (repo *MysqlTeamRepository) GetTeams(ctx context.Context, filter models.Team) (*[]models.Team, error) {
	var teams []models.Team
	if err := repo.db.WithContext(ctx).Where(filter).Find(&teams).Error; err != nil {
		return nil, fmt.Errorf("getTeams %+v: %v", filter, err.Error())
	}
	return &teams, nil
}

func (repo *MysqlTeamRepository) GetTeam(ctx context.Context, team *models.Team) error {
	return repo.db.WithContext(ctx).First(&team, &team.ID).Error
}

func (repo *MysqlTeamRepository) GetTeamByEntityID(ctx context.Context, entityId uint) (*models.Team, error) {
	var team models.Team
	if err := repo.db.WithContext(ctx).Model(&models.Team{}).
		Joins("LEFT JOIN members ON teams.id = members.team_id").
		Where("members.entity_id = ?", entityId).
		First(&team).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &team, nil
}

func (repo *MysqlTeamRepository) CreateTeam(ctx context.Context, team *models.Team, member *models.Member) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Create(team).Error; err != nil {
			return fmt.Errorf("createTeam1 %+v: %v", team, err)
		}
		member.TeamID = team.ID
		if err := tx.WithContext(ctx).Create(member).Error; err != nil {
			return fmt.Errorf("createTeam2 %+v: %v", team, err)
		}
		return nil
	})
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
