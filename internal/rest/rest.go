package rest

import (
	"github.com/codelix/ems/internal/rest/middleware"
	"github.com/codelix/ems/internal/rest/routes/entities"
	"github.com/codelix/ems/internal/rest/routes/invites"
	"github.com/codelix/ems/internal/rest/routes/members"
	"github.com/codelix/ems/internal/rest/routes/teams"
	"github.com/codelix/ems/internal/rest/routes/tokens"
	"github.com/codelix/ems/internal/service/entity"
	member "github.com/codelix/ems/internal/service/member"
	"github.com/codelix/ems/internal/service/team"
	"github.com/codelix/ems/internal/service/token"
	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	entityService entity.EntityService
	tokenService  token.TokenService
	teamService   team.TeamService
	memberService member.MemberService
}

func NewHttpServer(entityService entity.EntityService,
	tokenService token.TokenService,
	teamService team.TeamService,
	memberService member.MemberService) *HttpServer {
	return &HttpServer{entityService, tokenService, teamService, memberService}
}

func (server *HttpServer) Serve() {
	r := gin.Default()
	r.SetTrustedProxies(nil)
	r.Use(middleware.ErrorHandler())

	entitiesGroup := r.Group("entities")
	entities.Handle(entitiesGroup, &server.entityService, &server.memberService)
	tokenGroup := r.Group("token")
	tokens.Handle(tokenGroup, &server.tokenService)
	teamsGroup := r.Group("teams")
	teams.Handle(teamsGroup, &server.teamService, &server.memberService)
	invitesGroup := r.Group("invites")
	invites.Handle(invitesGroup, &server.memberService)
	membersGroup := r.Group("members")
	members.Handle(membersGroup, &server.memberService)

	r.Run("localhost:3006")
}
