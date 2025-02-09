package member

import (
	"context"

	"github.com/codelix/ems/pkg/models"
)

type MemberRepository interface {
	GetMembers(ctx context.Context) (*[]models.Member, error)
	GetMembersByTeamId(ctx context.Context, teamId uint) (*[]models.Member, error)
	CreateMember(ctx context.Context, member *models.Member) error
	DeleteMemberByEntityId(ctx context.Context, entityId uint) error
	GetMember(ctx context.Context, member *models.Member) error

	GetMemberInvites(ctx context.Context) (*[]models.MemberInvite, error)
	GetMemberInvite(ctx context.Context, inviteId uint) (*models.MemberInvite, error)
	GetMemberInvitesByInvitedId(ctx context.Context, invitedId uint) (*[]models.MemberInvite, error)
	CreateMemberInvite(ctx context.Context, memberInvite *models.MemberInvite) error
	DeclineMemberInvite(ctx context.Context, memberInvite models.MemberInvite) error
	AcceptMemberInvite(ctx context.Context, memberInvite models.MemberInvite) error
}
