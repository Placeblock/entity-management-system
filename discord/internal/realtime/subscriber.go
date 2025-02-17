package realtime

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Placeblock/nostalgicraft-discord/internal/service"
	"github.com/Placeblock/nostalgicraft-discord/pkg/config"
	perr "github.com/Placeblock/nostalgicraft-discord/pkg/errors"
	"github.com/Placeblock/nostalgicraft-discord/pkg/realtime"
	"github.com/Placeblock/nostalgicraft-ems/pkg/models"
	emsRealtime "github.com/Placeblock/nostalgicraft-ems/pkg/realtime"
	"github.com/bwmarrin/discordgo"
	colorful "github.com/lucasb-eyer/go-colorful"
	"github.com/pebbe/zmq4"
)

type Subscriber struct {
	cfg               *config.Config
	entityUserService *service.EntityUserService
	teamRoleService   *service.TeamRoleService
	discord           *discordgo.Session
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

	err = s.discord.GuildMemberNickname(s.cfg.Guild, userId, entity.Name)
	if err != nil {
		fmt.Print(fmt.Errorf("Could not nick user when renaming member: %v", err.Error()))
	}
}

func (s *Subscriber) onMemberCreate(member models.Member) {

}

func (s *Subscriber) onMemberRemove(member models.Member) {

}

func (s *Subscriber) onMemberInvite(invite models.MemberInvite) {

}

func (s *Subscriber) onMemberInviteAccept(invite models.MemberInvite) {

}

func (s *Subscriber) onMemberInviteDecline(invite models.MemberInvite) {

}

func (s *Subscriber) onTeamCreate(data emsRealtime.CreateTeamData) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	name := data.Team.Name
	params := discordgo.RoleParams{
		Name:  name,
		Color: getTeamColor(float64(*data.Team.Hue)),
	}
	role, err := s.discord.GuildRoleCreate(s.cfg.Guild, &params)
	if err != nil {
		fmt.Print(fmt.Errorf("Could not create Guild when creating team: %v", err.Error()))
		return
	}
	err = s.teamRoleService.CreateTeamRole(ctx, data.Team.ID, role.ID)
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
	err = s.discord.GuildMemberRoleAdd(s.cfg.Guild, userID, role.ID)
	if err != nil {
		fmt.Print(fmt.Errorf("Could not assign user to role when creating team: %v", err.Error()))
	}
}

func (s *Subscriber) onTeamRename(team models.Team) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	roleId, err := s.teamRoleService.GetRoleByTeamId(ctx, team.ID)
	if err != nil {
		fmt.Print(fmt.Errorf("Could not get Role ID when renaming team: %v", err.Error()))
		return
	}
	s.discord.GuildRoleEdit(s.cfg.Guild, roleId, &discordgo.RoleParams{
		Name: team.Name,
	})
}

func (s *Subscriber) onTeamRecolor(team models.Team) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	roleId, err := s.teamRoleService.GetRoleByTeamId(ctx, team.ID)
	if err != nil {
		fmt.Print(fmt.Errorf("Could not get Role ID when recoloring team: %v", err.Error()))
		return
	}
	s.discord.GuildRoleEdit(s.cfg.Guild, roleId, &discordgo.RoleParams{
		Color: getTeamColor(float64(*team.Hue)),
	})
}

func getTeamColor(hue float64) *int {
	c := colorful.Hcl(hue, 0.5, 0.5)
	r, g, b := c.RGB255()
	color := (int(r) << 16) + (int(g) << 8) + int(b)
	return &color
}
