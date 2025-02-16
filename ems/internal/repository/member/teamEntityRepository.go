package member

import (
	"context"

	"github.com/Placeblock/nostalgicraft-ems/pkg/models"
)

type MemberRepository interface {
	GetMembers(ctx context.Context, filter models.Member) (*[]models.Member, error)
	CreateMember(ctx context.Context, member *models.Member) error
	DeleteMember(ctx context.Context, member models.Member) error
	GetMember(ctx context.Context, member *models.Member) error

	GetMemberInvites(ctx context.Context, filter models.MemberInvite) (*[]models.MemberInvite, error)
	GetMemberInvite(ctx context.Context, invite *models.MemberInvite) error
	GetMemberInviteByInviterName(ctx context.Context, invite *models.MemberInvite, invitedId uint, inviterName string) error
	CreateMemberInvite(ctx context.Context, memberInvite *models.MemberInvite) error
	DeclineMemberInvite(ctx context.Context, memberInvite *models.MemberInvite) error
	AcceptMemberInvite(ctx context.Context, memberInvite *models.MemberInvite) (*models.Member, error)
}
