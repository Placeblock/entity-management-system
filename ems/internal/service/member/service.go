package member

import (
	"context"
	"fmt"

	"github.com/Placeblock/nostalgicraft-ems/internal/realtime"
	teamentity "github.com/Placeblock/nostalgicraft-ems/internal/repository/member"
	"github.com/Placeblock/nostalgicraft-ems/pkg/models"
	rtm "github.com/Placeblock/nostalgicraft-ems/pkg/realtime"
)

type MemberService struct {
	memberRepository *teamentity.MemberRepository
	publisher        *realtime.Publisher
}

func NewMemberService(repository teamentity.MemberRepository, publisher *realtime.Publisher) *MemberService {
	return &MemberService{&repository, publisher}
}

func (service *MemberService) SetTeam(ctx context.Context, entityId uint, teamId uint) error {
	member := models.Member{EntityID: entityId, TeamID: teamId}
	err := (*service.memberRepository).CreateMember(ctx, &member)
	if err != nil {
		return err
	}
	service.publisher.Channel <- rtm.Action{Type: "member.create", Data: member}
	return nil
}

func (service *MemberService) RemoveMember(ctx context.Context, member *models.Member) error {
	err := (*service.memberRepository).GetMember(ctx, member)
	if err != nil {
		return err
	}
	err = (*service.memberRepository).DeleteMember(ctx, *member)
	if err != nil {
		return err
	}
	service.publisher.Channel <- rtm.Action{Type: "member.remove", Data: member}
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

func (service *MemberService) GetMemberByEntityId(ctx context.Context, entityId uint) (*models.Member, error) {
	member := models.Member{EntityID: entityId}
	err := (*service.memberRepository).GetMember(ctx, &member)
	if err != nil {
		return nil, err
	}
	if member.ID == 0 {
		return nil, nil
	}
	return &member, nil
}

func (service *MemberService) GetEntityIdByMemberId(ctx context.Context, memberId uint) (*uint, error) {
	member := models.Member{ID: memberId}
	err := (*service.memberRepository).GetMember(ctx, &member)
	if err != nil {
		return nil, err
	}
	return &member.EntityID, nil
}

func (service *MemberService) CreateInvite(ctx context.Context, invitedId uint, inviterId uint) (*models.MemberInvite, error) {
	memberInvite := models.MemberInvite{InvitedID: invitedId, InviterID: inviterId}
	err := (*service.memberRepository).CreateMemberInvite(ctx, &memberInvite)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%+v\n", memberInvite)
	service.publisher.Channel <- rtm.Action{Type: "member.invite", Data: memberInvite}
	return &memberInvite, nil
}

func (service *MemberService) AcceptInvite(ctx context.Context, inviteId uint) (*models.Member, error) {
	memberInvite := models.MemberInvite{ID: inviteId}
	member, err := (*service.memberRepository).AcceptMemberInvite(ctx, &memberInvite)
	if err != nil {
		return nil, err
	}
	service.publisher.Channel <- rtm.Action{Type: "member.invite.accept", Data: memberInvite}
	service.publisher.Channel <- rtm.Action{Type: "member.create", Data: member}
	return member, nil
}

func (service *MemberService) DeclineInvite(ctx context.Context, inviteId uint) error {
	memberInvite := models.MemberInvite{ID: inviteId}
	err := (*service.memberRepository).DeclineMemberInvite(ctx, &memberInvite)
	if err != nil {
		return err
	}
	service.publisher.Channel <- rtm.Action{Type: "member.invite.decline", Data: memberInvite}
	return nil
}

func (service *MemberService) GetMemberInvitesByInvitedId(ctx context.Context, invitedId uint) (*[]models.MemberInvite, error) {
	return (*service.memberRepository).GetMemberInvites(ctx, models.MemberInvite{InvitedID: invitedId})
}

func (service *MemberService) GetMemberInvite(ctx context.Context, inviteId uint) (*models.MemberInvite, error) {
	invite := models.MemberInvite{ID: inviteId}
	err := (*service.memberRepository).GetMemberInvite(ctx, &invite)
	if err != nil {
		return nil, err
	}
	return &invite, nil
}

func (service *MemberService) GetMemberInviteByInviterName(ctx context.Context, invitedId uint, inviterName string) (*models.MemberInvite, error) {
	var invite models.MemberInvite
	err := (*service.memberRepository).GetMemberInviteByInviterName(ctx, &invite, invitedId, inviterName)
	if err != nil {
		return nil, err
	}
	if invite.ID == 0 {
		return nil, nil
	}
	return &invite, nil
}

func (service *MemberService) CreateMessage(ctx context.Context, memberId uint, message string) error {
	teamMessage := models.TeamMessage{MemberID: memberId, Message: message}
	err := (*service.memberRepository).CreateMessage(ctx, &teamMessage)
	if err != nil {
		return err
	}
	service.publisher.Channel <- rtm.Action{Type: "team.message", Data: teamMessage}
	return nil
}
