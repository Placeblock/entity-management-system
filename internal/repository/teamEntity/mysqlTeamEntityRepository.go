package teamentity

import (
	"context"
	"fmt"

	"github.com/codelix/ems/pkg/models"
	"gorm.io/gorm"
)

type MysqlTeamEntityRepository struct {
	db *gorm.DB
}

func NewMysqlTeamEntityRepository(db *gorm.DB) *MysqlTeamEntityRepository {
	return &MysqlTeamEntityRepository{db}
}

func (repo *MysqlTeamEntityRepository) GetTeamEntities(ctx context.Context) (*[]models.TeamEntity, error) {
	var teamEntities []models.TeamEntity
	if err := repo.db.WithContext(ctx).Preload("Team").Preload("Entity").Find(&teamEntities).Error; err != nil {
		return nil, fmt.Errorf("getTeamEntities: %v", err.Error())
	}
	return &teamEntities, nil
}

func (repo *MysqlTeamEntityRepository) GetTeamEntitiesByTeamId(ctx context.Context, teamId uint) (*[]models.TeamEntity, error) {
	var teamEntities []models.TeamEntity
	if err := repo.db.WithContext(ctx).Preload("Team").Preload("Entity").Find(&teamEntities).Error; err != nil {
		return nil, fmt.Errorf("getTeamEntities: %v", err.Error())
	}
	return &teamEntities, nil
}

func (repo *MysqlTeamEntityRepository) GetTeamEntityByEntityId(ctx context.Context, entityId uint) (*models.TeamEntity, error) {
	var teamEntity models.TeamEntity
	if err := repo.db.WithContext(ctx).First(&teamEntity, "EntityID = ?", entityId).Error; err != nil {
		return nil, fmt.Errorf("getTeamEntityByEntityId %d: %s", entityId, err.Error())
	}
	return &teamEntity, nil
}

func (repo *MysqlTeamEntityRepository) CreateTeamEntity(ctx context.Context, teamEntity *models.TeamEntity) error {
	if err := repo.db.WithContext(ctx).Create(teamEntity).Error; err != nil {
		return fmt.Errorf("createEntity: %v", err.Error())
	}
	return nil
}

func (repo *MysqlTeamEntityRepository) DeleteTeamEntity(ctx context.Context, entityId uint) error {
	if err := repo.db.WithContext(ctx).Delete(models.TeamEntity{EntityID: entityId}).Error; err != nil {
		return fmt.Errorf("deleteEntity %d: %v", entityId, err.Error())
	}
	return nil
}
