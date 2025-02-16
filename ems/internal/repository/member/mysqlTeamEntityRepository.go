package member

import (
	"context"
	"errors"
	"fmt"

	"github.com/Placeblock/nostalgicraft-ems/pkg/models"
	"gorm.io/gorm"
)

type MysqlMemberRepository struct {
	db *gorm.DB
}

func NewMysqlMemberRepository(db *gorm.DB) *MysqlMemberRepository {
	return &MysqlMemberRepository{db}
}

func (repo *MysqlMemberRepository) GetMembers(ctx context.Context, filter models.Member) (*[]models.Member, error) {
	var teamEntities []models.Member
	if err := repo.db.WithContext(ctx).Where(filter).Find(&teamEntities).Error; err != nil {
		return nil, fmt.Errorf("getMembers: %v", err.Error())
	}
	return &teamEntities, nil
}

func (repo *MysqlMemberRepository) GetMember(ctx context.Context, member *models.Member) error {
	if err := repo.db.WithContext(ctx).First(&member).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
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
	return repo.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Find(&member).Error; err != nil {
			return fmt.Errorf("deleteMember1: %v", err.Error())
		}
		fmt.Printf("%+v\n", member)
		if err := tx.WithContext(ctx).Delete(member).Error; err != nil {
			return fmt.Errorf("deleteMember2: %v", err.Error())
		}
		if err := tx.WithContext(ctx).Where(&models.MemberInvite{InviterID: member.EntityID}).Delete(&models.MemberInvite{}).Error; err != nil {
			return fmt.Errorf("deleteMember3: %v", err.Error())
		}
		countMember := models.Member{TeamID: member.TeamID}
		var memberCount int64
		if err := tx.WithContext(ctx).Model(&models.Member{}).Where(&countMember).Count(&memberCount).Error; err != nil {
			return fmt.Errorf("deleteMember4: %v", err.Error())
		}
		if memberCount != 0 {
			return nil
		}
		if err := tx.WithContext(ctx).Delete(models.Team{ID: member.TeamID}).Error; err != nil {
			return fmt.Errorf("deleteMember5: %v", err.Error())
		}
		return nil
	})
}

func (repo *MysqlMemberRepository) GetMemberInvites(ctx context.Context, filter models.MemberInvite) (*[]models.MemberInvite, error) {
	var memberInvites []models.MemberInvite
	if err := repo.db.WithContext(ctx).Where(filter).Find(&memberInvites).Error; err != nil {
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
	if err := repo.db.WithContext(ctx).Where(invite).Delete(invite).Error; err != nil {
		return fmt.Errorf("declineMemberInvite %+v: %v", invite, err.Error())
	}
	return nil
}

func (repo *MysqlMemberRepository) AcceptMemberInvite(ctx context.Context, invite models.MemberInvite) (*models.Member, error) {
	var member models.Member
	err := repo.db.Transaction(func(tx *gorm.DB) error {
		inviter := models.Member{EntityID: invite.InviterID}
		if err := tx.WithContext(ctx).First(&inviter).Error; err != nil {
			return fmt.Errorf("acceptMemberInvite1 %+v: %v", invite, err.Error())
		}
		if err := tx.WithContext(ctx).Delete(invite).Error; err != nil {
			return fmt.Errorf("acceptMemberInvite2 %+v: %v", invite, err.Error())
		}
		member = models.Member{EntityID: invite.InvitedID, TeamID: inviter.TeamID}
		if err := tx.WithContext(ctx).Create(&member).Error; err != nil {
			return fmt.Errorf("acceptMemberInvite3 %+v: %v", invite, err.Error())
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &member, nil
}
