package member

import (
	"context"
	"fmt"

	"github.com/codelix/ems/pkg/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MysqlMemberRepository struct {
	db *gorm.DB
}

func NewMysqlMemberRepository(db *gorm.DB) *MysqlMemberRepository {
	return &MysqlMemberRepository{db}
}

func (repo *MysqlMemberRepository) GetMembers(ctx context.Context, filter models.Member) (*[]models.Member, error) {
	var teamEntities []models.Member
	if err := repo.db.WithContext(ctx).Where(filter).Preload(clause.Associations).Find(&teamEntities).Error; err != nil {
		return nil, fmt.Errorf("getMembers: %v", err.Error())
	}
	return &teamEntities, nil
}

func (repo *MysqlMemberRepository) GetMember(ctx context.Context, member *models.Member) error {
	if err := repo.db.WithContext(ctx).First(&member).Error; err != nil {
		return fmt.Errorf("getMember: %s", err.Error())
	}
	return nil
}

func (repo *MysqlMemberRepository) CreateMember(ctx context.Context, member *models.Member) error {
	if err := repo.db.WithContext(ctx).Create(member).Error; err != nil {
		return fmt.Errorf("createMember: %v", err.Error())
	}
	return nil
}

func (repo *MysqlMemberRepository) DeleteMember(ctx context.Context, member models.Member) error {
	if err := repo.db.WithContext(ctx).Delete(member).Error; err != nil {
		return fmt.Errorf("deleteMember: %v", err.Error())
	}
	return nil
}

func (repo *MysqlMemberRepository) GetMemberInvites(ctx context.Context, filter models.MemberInvite) (*[]models.MemberInvite, error) {
	var memberInvites []models.MemberInvite
	if err := repo.db.WithContext(ctx).Where(filter).Preload(clause.Associations).Find(&memberInvites).Error; err != nil {
		return nil, fmt.Errorf("getMemberInvites %+v: %v", filter, err.Error())
	}
	return &memberInvites, nil
}

func (repo *MysqlMemberRepository) GetMemberInvite(ctx context.Context, invite *models.MemberInvite) error {
	if err := repo.db.WithContext(ctx).First(invite).Error; err != nil {
		return fmt.Errorf("getMemberInvite %+v: %v", invite, err.Error())
	}
	return nil
}

func (repo *MysqlMemberRepository) CreateMemberInvite(ctx context.Context, invite *models.MemberInvite) error {
	if err := repo.db.WithContext(ctx).Create(invite).Error; err != nil {
		return fmt.Errorf("createMemberInvite %+v: %v", invite, err.Error())
	}
	return nil
}

func (repo *MysqlMemberRepository) DeclineMemberInvite(ctx context.Context, invite models.MemberInvite) error {
	if err := repo.db.WithContext(ctx).Delete(invite).Error; err != nil {
		return fmt.Errorf("declineMemberInvite %+v: %v", invite, err.Error())
	}
	return nil
}

func (repo *MysqlMemberRepository) AcceptMemberInvite(ctx context.Context, invite models.MemberInvite) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Delete(invite).Error; err != nil {
			return fmt.Errorf("acceptMemberInvite %+v: %v", invite, err.Error())
		}
		member := models.Member{EntityID: invite.InvitedID, TeamID: invite.TeamID}
		if err := tx.WithContext(ctx).Create(member).Error; err != nil {
			return fmt.Errorf("acceptMemberInvite %+v: %v", invite, err.Error())
		}
		return nil
	})
}
