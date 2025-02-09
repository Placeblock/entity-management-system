package teamentity

import (
	"context"
	"fmt"

	"github.com/codelix/ems/pkg/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MysqlTeamEntityRepository struct {
	db *gorm.DB
}

func NewMysqlTeamEntityRepository(db *gorm.DB) *MysqlTeamEntityRepository {
	return &MysqlTeamEntityRepository{db}
}

func (repo *MysqlTeamEntityRepository) GetTeamEntities(ctx context.Context) (*[]models.TeamEntity, error) {
	var teamEntities []models.TeamEntity
	if err := repo.db.WithContext(ctx).Preload(clause.Associations).Find(&teamEntities).Error; err != nil {
		return nil, fmt.Errorf("getTeamEntities: %v", err.Error())
	}
	return &teamEntities, nil
}

func (repo *MysqlTeamEntityRepository) GetTeamEntitiesByTeamId(ctx context.Context, teamId uint) (*[]models.TeamEntity, error) {
	var teamEntities []models.TeamEntity
	if err := repo.db.WithContext(ctx).Preload(clause.Associations).Find(&teamEntities).Error; err != nil {
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
		return fmt.Errorf("createTeamEntity: %v", err.Error())
	}
	return nil
}

func (repo *MysqlTeamEntityRepository) DeleteTeamEntity(ctx context.Context, entityId uint) error {
	if err := repo.db.WithContext(ctx).Delete(models.TeamEntity{EntityID: entityId}).Error; err != nil {
		return fmt.Errorf("deleteTeamEntity %d: %v", entityId, err.Error())
	}
	return nil
}

func (repo *MysqlTeamEntityRepository) GetTeamEntityInvites(ctx context.Context) (*[]models.TeamEntityInvite, error) {
	var teamEntityInvites []models.TeamEntityInvite
	if err := repo.db.WithContext(ctx).Preload(clause.Associations).Find(&teamEntityInvites).Error; err != nil {
		return nil, fmt.Errorf("getTeamEntityInvites: %v", err.Error())
	}
	return &teamEntityInvites, nil
}

func (repo *MysqlTeamEntityRepository) GetTeamEntityInvitesByInvitedId(ctx context.Context, invitedId uint) (*[]models.TeamEntityInvite, error) {
	var teamEntityInvites []models.TeamEntityInvite
	if err := repo.db.WithContext(ctx).Where(models.TeamEntityInvite{InvitedID: invitedId}).Find(&teamEntityInvites).Error; err != nil {
		return nil, fmt.Errorf("getTeamEntityInvitesByInvitedId %d: %s", invitedId, err.Error())
	}
	return &teamEntityInvites, nil
}

func (repo *MysqlTeamEntityRepository) GetTeamEntityInvite(ctx context.Context, invitedId uint, teamId uint) (*models.TeamEntityInvite, error) {
	var teamEntityInvite models.TeamEntityInvite
	if err := repo.db.WithContext(ctx).Where(models.TeamEntityInvite{InvitedID: invitedId, TeamID: teamId}).First(&teamEntityInvite).Error; err != nil {
		return nil, fmt.Errorf("getTeamEntityInvite %d: %s", invitedId, err.Error())
	}
	return &teamEntityInvite, nil
}

func (repo *MysqlTeamEntityRepository) CreateTeamEntityInvite(ctx context.Context, invite models.TeamEntityInvite) error {
	if err := repo.db.WithContext(ctx).Create(invite).Error; err != nil {
		return fmt.Errorf("createTeamEntityInvite: %v", err.Error())
	}
	return nil
}

func (repo *MysqlTeamEntityRepository) DeclineTeamEntityInvite(ctx context.Context, teamEntityInvite models.TeamEntityInvite) error {
	if err := repo.db.WithContext(ctx).Delete(teamEntityInvite).Error; err != nil {
		return fmt.Errorf("declineTeamEntityInvite %d: %v", teamEntityInvite.InvitedID, err.Error())
	}
	return nil
}

func (repo *MysqlTeamEntityRepository) AcceptTeamEntityInvite(ctx context.Context, teamEntityInvite models.TeamEntityInvite) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Delete(teamEntityInvite).Error; err != nil {
			return fmt.Errorf("acceptTeamEntityInvite %d: %v", teamEntityInvite.InvitedID, err.Error())
		}
		teamEntity := models.TeamEntity{EntityID: teamEntityInvite.InvitedID, TeamID: teamEntityInvite.TeamID}
		if err := tx.WithContext(ctx).Create(teamEntity).Error; err != nil {
			return fmt.Errorf("acceptTeamEntityInvite: %v", err.Error())
		}
		return nil
	})
}
