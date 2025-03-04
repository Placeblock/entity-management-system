package realtime

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Placeblock/nostalgicraft-discord/internal"
	"github.com/Placeblock/nostalgicraft-discord/internal/service"
	perr "github.com/Placeblock/nostalgicraft-discord/pkg/errors"
	"github.com/Placeblock/nostalgicraft-discord/pkg/realtime"
	"github.com/Placeblock/nostalgicraft-ems/pkg/models"
	emsRealtime "github.com/Placeblock/nostalgicraft-ems/pkg/realtime"
	"github.com/bwmarrin/discordgo"
	colorful "github.com/lucasb-eyer/go-colorful"
	"github.com/pebbe/zmq4"
)

type Subscriber struct {
	entityUserService *service.EntityUserService
	teamDataService   *service.TeamDataService
	discord           *discordgo.Session
}

func NewSubscriber(entityUserService *service.EntityUserService,
	teamDataService *service.TeamDataService, discord *discordgo.Session) *Subscriber {
	return &Subscriber{entityUserService, teamDataService, discord}
}

func (s *Subscriber) Listen() {
	zctx, err := zmq4.NewContext()
	socket, err := zctx.NewSocket(zmq4.SUB)
	if err != nil {
		log.Fatal("Could not create ZMQ Socket")
		return
	}
	err = socket.Connect("tcp://127.0.0.1:3008")
	if err != nil {
		log.Fatal("Could not connect to ZMQ Publisher: ", err)
		return
	}
	err = socket.SetSubscribe("")
	if err != nil {
		log.Fatal("Could not set Subscribe Filter: ", err)
		return
	}
	fmt.Println("Connected to ZMQ Socket!")
	for {
		bytes, err := socket.RecvBytes(0)
		fmt.Printf("%s\n", string(bytes[:]))
		if err != nil {
			fmt.Printf("Could not recv data: %v\n", err)
			continue
		}
		var msg realtime.RawAction
		err = json.Unmarshal(bytes, &msg)
		if err != nil {
			fmt.Printf("Could not unmarshal data: %v\n", err)
			continue
		}
		switch msg.Type {
		case "entity.rename":
			var entity models.Entity
			err = json.Unmarshal(msg.Data, &entity)
			s.onEntityRename(entity)
		case "member.create":
			var member models.Member
			err = json.Unmarshal(msg.Data, &member)
			s.onMemberCreate(member)
		case "member.remove":
			var member models.Member
			err = json.Unmarshal(msg.Data, &member)
			s.onMemberRemove(member)
		case "member.invite":
			var invite models.MemberInvite
			err = json.Unmarshal(msg.Data, &invite)
			s.onMemberInvite(invite)
		case "member.invite.accept":
			var invite models.MemberInvite
			err = json.Unmarshal(msg.Data, &invite)
			s.onMemberInviteAccept(invite)
		case "member.invite.decline":
			var invite models.MemberInvite
			err = json.Unmarshal(msg.Data, &invite)
			s.onMemberInviteDecline(invite)
		case "team.create":
			var teamCreateData emsRealtime.CreateTeamData
			err = json.Unmarshal(msg.Data, &teamCreateData)
			s.onTeamCreate(teamCreateData)
		case "team.rename":
			var team models.Team
			err = json.Unmarshal(msg.Data, &team)
			s.onTeamRename(team)
		case "team.recolor":
			var team models.Team
			err = json.Unmarshal(msg.Data, &team)
			s.onTeamRecolor(team)
		case "team.message":
			var teamMessage models.TeamMessage
			err = json.Unmarshal(msg.Data, &teamMessage)
			s.onTeamMessage(teamMessage)
		}
	}
}

func (s *Subscriber) onEntityRename(entity models.Entity) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userId, err := s.entityUserService.GetUserIdByEntityId(ctx, entity.ID)
	if err != nil {
		if errors.Is(err, perr.ErrNotFound{}) {
			return
		}
		fmt.Print(fmt.Errorf("Could not get User ID when renaming entity: %v", err.Error()))
		return
	}

	err = s.discord.GuildMemberNickname(internal.Config.Guild, userId, entity.Name)
	if err != nil {
		fmt.Print(fmt.Errorf("Could not nick user when renaming member: %v", err.Error()))
	}
}

func (s *Subscriber) onMemberCreate(member models.Member) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userId, err := s.entityUserService.GetUserIdByEntityId(ctx, member.EntityID)
	if err != nil {
		if errors.Is(err, perr.ErrNotFound{}) {
			return
		}
		fmt.Print(fmt.Errorf("Could not get User ID when creating team member: %v", err.Error()))
		return
	}

	teamData, err := s.teamDataService.GetTeamDataByTeamId(ctx, member.TeamID)
	if err != nil {
		fmt.Print(fmt.Errorf("Could not get Role ID when creating team member: %v", err.Error()))
		return
	}

	err = s.discord.GuildMemberRoleAdd(internal.Config.Guild, userId, teamData.RoleID)
	if err != nil {
		fmt.Print(fmt.Errorf("Could not add User to Role when creating team member: %v", err.Error()))
	}
}

func (s *Subscriber) onMemberRemove(member models.Member) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userId, err := s.entityUserService.GetUserIdByEntityId(ctx, member.EntityID)
	if err != nil {
		if errors.Is(err, perr.ErrNotFound{}) {
			return
		}
		fmt.Print(fmt.Errorf("Could not get User ID when removing team member: %v", err.Error()))
		return
	}

	teamData, err := s.teamDataService.GetTeamDataByTeamId(ctx, member.TeamID)
	if err != nil {
		fmt.Print(fmt.Errorf("Could not get Role ID when removing team member: %v", err.Error()))
		return
	}

	err = s.discord.GuildMemberRoleRemove(internal.Config.Guild, userId, teamData.RoleID)
	if err != nil {
		fmt.Print(fmt.Errorf("Could not remove User from Role when removing team member: %v", err.Error()))
	}
}

func (s *Subscriber) onMemberInvite(invite models.MemberInvite) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	invitedUserId, err := s.entityUserService.GetUserIdByEntityId(ctx, invite.InvitedID)
	if err != nil {
		if errors.Is(err, perr.ErrNotFound{}) {
			return
		}
		fmt.Print(fmt.Errorf("Could not get Invited User ID when receiving invite: %v", err.Error()))
		return
	}

	userChan, err := s.discord.UserChannelCreate(invitedUserId)
	if err != nil {
		fmt.Print(fmt.Errorf("Could not create Private Channel when receiving invite: %v", err.Error()))
		return
	}
	inviterName := invite.Inviter.Entity.Name
	serializedInviteId := strconv.FormatUint(uint64(invite.ID), 10)
	components := []discordgo.MessageComponent{
		discordgo.Button{
			Label:    "ACCEPT",
			Style:    discordgo.PrimaryButton,
			CustomID: "accept-invite-" + serializedInviteId,
		},
		discordgo.Button{
			Label:    "DECLINE",
			Style:    discordgo.DangerButton,
			CustomID: "decline-invite-" + serializedInviteId,
		},
	}
	_, err = s.discord.ChannelMessageSendComplex(userChan.ID, &discordgo.MessageSend{
		Content: "```ansi\n\u001b[1;35m" + inviterName + " \u001b[0minvited you to their Team```",
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: components,
			},
		},
	})
	if err != nil {
		fmt.Print(fmt.Errorf("Could not send Message when receiving invite: %v", err.Error()))
		return
	}
}

func (s *Subscriber) onMemberInviteAccept(invite models.MemberInvite) {
	s.sendMemberInviteProcessMessage(invite, true)
}

func (s *Subscriber) onMemberInviteDecline(invite models.MemberInvite) {
	s.sendMemberInviteProcessMessage(invite, false)
}

func (s *Subscriber) sendMemberInviteProcessMessage(invite models.MemberInvite, accepted bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	invitedName := invite.Invited.Name
	inviterName := invite.Inviter.Entity.Name
	var message string
	if accepted {
		message = "\u001b[1;32maccepted"
	} else {
		message = "\u001b[1;31mdeclined"
	}

	invitedUserId, err := s.entityUserService.GetUserIdByEntityId(ctx, invite.Invited.ID)
	if err == nil {
		userChan, err := s.discord.UserChannelCreate(invitedUserId)
		if err == nil {
			_, err = s.discord.ChannelMessageSend(userChan.ID,
				"```ansi\nYou "+message+" \u001b[1;35m"+inviterName+"'s \u001b[0mTeam-Invite.```")
			if err != nil {
				fmt.Print(fmt.Errorf("Could not send invited message when processing invite: %v", err.Error()))
			}
		} else {
			fmt.Print(fmt.Errorf("Could not create invited channel when processing invite: %v", err.Error()))
		}
	}
	inviterUserId, err := s.entityUserService.GetUserIdByEntityId(ctx, invite.Inviter.EntityID)
	if err == nil {
		userChan, err := s.discord.UserChannelCreate(inviterUserId)
		if err == nil {
			_, err = s.discord.ChannelMessageSend(userChan.ID,
				"```ansi\n\u001b[1;35m"+invitedName+" "+message+" \u001b[0myour Team-Invite.```")
			if err != nil {
				fmt.Print(fmt.Errorf("Could not send inviter message when processing invite: %v", err.Error()))
			}
		} else {
			fmt.Print(fmt.Errorf("Could not create inviter channel when processing invite: %v", err.Error()))
		}
	}
}

func (s *Subscriber) onTeamCreate(data emsRealtime.CreateTeamData) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	name := data.Team.Name
	hoist := true
	params := discordgo.RoleParams{
		Name:  name,
		Color: getTeamColor(float64(*data.Team.Hue)),
		Hoist: &hoist,
	}
	role, err := s.discord.GuildRoleCreate(internal.Config.Guild, &params)
	if err != nil {
		fmt.Print(fmt.Errorf("Could not create Role when creating team: %v", err.Error()))
		return
	}
	channel, err := s.discord.GuildChannelCreateComplex(internal.Config.Guild, discordgo.GuildChannelCreateData{
		Name:     name,
		Type:     discordgo.ChannelTypeGuildText,
		ParentID: internal.Config.TeamsCategoryID,
		PermissionOverwrites: []*discordgo.PermissionOverwrite{
			{
				ID:    role.ID,
				Type:  discordgo.PermissionOverwriteTypeRole,
				Allow: discordgo.PermissionViewChannel,
				Deny:  0,
			},
			{
				ID:   internal.Config.EveryoneRoleID,
				Deny: discordgo.PermissionViewChannel,
			},
		},
	})
	if err != nil {
		fmt.Print(fmt.Errorf("Could not create Guild when creating team: %v", err.Error()))
		return
	}
	err = s.teamDataService.CreateTeamData(ctx, data.Team.ID, role.ID, channel.ID)
	if err != nil {
		fmt.Print(fmt.Errorf("Could not create TeamRole when creating team: %v", err.Error()))
		return
	}
	userID, err := s.entityUserService.GetUserIdByEntityId(ctx, data.Member.EntityID)
	if err != nil {
		if errors.Is(err, perr.ErrNotFound{}) {
			return
		}
		fmt.Print(fmt.Errorf("Could not get user id for entity id when creating team: %v", err.Error()))
		return
	}
	err = s.discord.GuildMemberRoleAdd(internal.Config.Guild, userID, role.ID)
	if err != nil {
		fmt.Print(fmt.Errorf("Could not assign user to role when creating team: %v", err.Error()))
	}
}

func (s *Subscriber) onTeamRename(team models.Team) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	teamData, err := s.teamDataService.GetTeamDataByTeamId(ctx, team.ID)
	if err != nil {
		fmt.Print(fmt.Errorf("Could not get Role ID when renaming team: %v", err.Error()))
		return
	}
	_, err = s.discord.GuildRoleEdit(internal.Config.Guild, teamData.RoleID, &discordgo.RoleParams{
		Name: team.Name,
	})
	if err != nil {
		fmt.Print(fmt.Errorf("Could not edit Role when renaming team: %v", err.Error()))
		return
	}
	_, err = s.discord.ChannelEdit(teamData.ChannelID, &discordgo.ChannelEdit{
		Name: team.Name,
	})
	if err != nil {
		fmt.Print(fmt.Errorf("Could not edit Channel when renaming team: %v", err.Error()))
		return
	}
}

func (s *Subscriber) onTeamRecolor(team models.Team) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	teamData, err := s.teamDataService.GetTeamDataByTeamId(ctx, team.ID)
	if err != nil {
		fmt.Print(fmt.Errorf("Could not get Role ID when recoloring team: %v", err.Error()))
		return
	}
	_, err = s.discord.GuildRoleEdit(internal.Config.Guild, teamData.RoleID, &discordgo.RoleParams{
		Color: getTeamColor(float64(*team.Hue)),
	})
	if err != nil {
		fmt.Print(fmt.Errorf("Could not edit Role when recoloring team: %v", err.Error()))
	}
}

func (s *Subscriber) onTeamMessage(teamMessage models.TeamMessage) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	teamData, err := s.teamDataService.GetTeamDataByTeamId(ctx, teamMessage.Member.TeamID)
	if err != nil {
		fmt.Print(fmt.Errorf("Could not get Team Data when receiving team message: %v", err.Error()))
		return
	}
	senderName := teamMessage.Member.Entity.Name
	_, err = s.discord.ChannelMessageSend(teamData.ChannelID, "**"+senderName+":**  "+teamMessage.Message)
	if err != nil {
		fmt.Print(fmt.Errorf("Could not send Team Message when receiving team message: %v", err.Error()))
	}
}

func getTeamColor(hue float64) *int {
	c := colorful.Hsv(hue*360, 0.7, 0.5)
	r, g, b := c.RGB255()
	color := (int(r) << 16) + (int(g) << 8) + int(b)
	return &color
}
