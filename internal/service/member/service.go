package member

import (
	"context"

	"github.com/codelix/ems/internal/realtime"
	teamentity "github.com/codelix/ems/internal/repository/member"
	"github.com/codelix/ems/pkg/models"
)

type MemberService struct {
	memberRepository *teamentity.MemberRepository
	publisher        *realtime.Publisher
}

func NewMemberService(repository teamentity.MemberRepository, publisher *realtime.Publisher) *MemberService {
	return &MemberService{&repository, publisher}
}

func (service *MemberService) SetTeam(ctx context.Context, entityId uint, teamId uint) error {
	teamEntity := models.Member{EntityID: entityId, TeamID: teamId}
	err := (*service.memberRepository).CreateMember(ctx, &teamEntity)
	if err != nil {
		return err
	}
	service.publisher.Channel <- realtime.Action{Type: "member.create", Data: teamEntity}
	return nil
}

func (service *MemberService) RemoveMemberByMemberId(ctx context.Context, memberId uint) error {
	member := models.Member{ID: memberId}
	err := (*service.memberRepository).GetMember(ctx, &member)
	if err != nil {
		return err
	}
	err = (*service.memberRepository).DeleteMember(ctx, member)
	if err != nil {
		return err
	}
	service.publisher.Channel <- realtime.Action{Type: "member.remove", Data: member}
	return nil
}

func (service *MemberService) LeaveTeam(ctx context.Context, entityId uint) error {
	member := models.Member{EntityID: entityId}
	err := (*service.memberRepository).GetMember(ctx, &member)
	if err != nil {
		return err
	}
	err = (*service.memberRepository).DeleteMember(ctx, member)
	if err != nil {
		return err
	}
	service.publisher.Channel <- realtime.Action{Type: "member.remove", Data: member}
	return nil
}

func (service *MemberService) GetMembers(ctx context.Context) (*[]models.Member, error) {
	return (*service.memberRepository).GetMembers(ctx, models.Member{})
}

func (service *MemberService) GetMembersByTeamId(ctx context.Context, teamId uint) (*[]models.Member, error) {
	return (*service.memberRepository).GetMembers(ctx, models.Member{TeamID: teamId})
}

func (service *MemberService) GetMember(ctx context.Context, memberId uint) (*models.Member, error) {
	member := models.Member{ID: memberId}
	err := (*service.memberRepository).GetMember(ctx, &member)
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (service *MemberService) CreateInvite(ctx context.Context, invitedId uint, inviterId uint, teamId uint) (*models.MemberInvite, error) {
	memberInvite := models.MemberInvite{InvitedID: invitedId, InviterID: inviterId, TeamID: teamId}
	err := (*service.memberRepository).CreateMemberInvite(ctx, &memberInvite)
	if err != nil {
		return nil, err
	}
	return &memberInvite, nil
}

func (service *MemberService) ProcessInvite(ctx context.Context, inviteId uint, accept bool) error {
	memberInvite := models.MemberInvite{ID: inviteId}
	if accept {
		return (*service.memberRepository).AcceptMemberInvite(ctx, memberInvite)
	} else {
		return (*service.memberRepository).DeclineMemberInvite(ctx, memberInvite)
	}
}

func (service *MemberService) GetMemberInvitesByInvitedId(ctx context.Context, invitedId uint) (*[]models.MemberInvite, error) {
	return (*service.memberRepository).GetMemberInvites(ctx, models.MemberInvite{InvitedID: invitedId})
}

func (service *MemberService) GetMemberInvitesByTeamId(ctx context.Context, teamId uint) (*[]models.MemberInvite, error) {
	return (*service.memberRepository).GetMemberInvites(ctx, models.MemberInvite{TeamID: teamId})
}
