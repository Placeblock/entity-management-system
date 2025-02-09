package teams

import (
	"context"
	"net/http"
	"strconv"
	"time"

	member "github.com/codelix/ems/internal/service/member"
	"github.com/codelix/ems/internal/service/team"
	"github.com/codelix/ems/pkg/models"
	"github.com/codelix/ems/pkg/rest"
	"github.com/gin-gonic/gin"
)

func Handle(g *gin.RouterGroup, teamService *team.TeamService, memberService *member.MemberService) {
	g.GET("", func(ctx *gin.Context) {
		getTeams(ctx, teamService)
	})
	g.POST("", func(ctx *gin.Context) {
		createTeam(ctx, teamService)
	})
	g.GET(":id", func(ctx *gin.Context) {
		getTeam(ctx, teamService)
	})
	g.PUT(":id/owner", func(ctx *gin.Context) {
		setOwner(ctx, teamService)
	})
	g.PUT(":id/color", func(ctx *gin.Context) {
		recolorTeam(ctx, teamService)
	})
	g.PUT(":id/name", func(ctx *gin.Context) {
		renameTeam(ctx, teamService)
	})
	g.GET(":id/members", func(ctx *gin.Context) {
		getMembers(ctx, memberService)
	})

	g.GET(":id/invites", func(ctx *gin.Context) {
		getInvites(ctx, memberService)
	})
	g.POST(":id/invites", func(ctx *gin.Context) {
		inviteEntity(ctx, memberService)
	})
}

type createParams struct {
	Name    string      `json:"name" binding:"required"`
	Hue     *models.Hue `json:"hue" binding:"required"`
	OwnerID uint        `json:"owner_id" binding:"required"`
}

func createTeam(ctx *gin.Context, teamService *team.TeamService) {
	var params createParams
	err := ctx.ShouldBindJSON(&params)
	if err != nil {
		ctx.Error(&rest.HTTPError{Title: "Invalid Parameters", Detail: "No or invalid parameters where provided to create the team", Status: http.StatusBadRequest, Cause: err})
		return
	}
	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	team := models.Team{Name: params.Name, Hue: params.Hue, OwnerID: params.OwnerID}
	err = teamService.CreateTeam(context, &team)
	cancel()
	if err != nil {
		ctx.Error(&rest.HTTPError{Title: "Unexpected Error", Detail: "An unexpected Error occurde while creating the team", Status: http.StatusInternalServerError, Cause: err})
		return
	}
	ctx.JSON(http.StatusOK, rest.Response{Data: team})
}

type renameParams struct {
	Name string `json:"name" binding:"required"`
}

func renameTeam(ctx *gin.Context, teamService *team.TeamService) {
	serializedId := ctx.Param("id")
	id, err := strconv.ParseUint(serializedId, 10, 0)
	if err != nil {
		ctx.Error(&rest.HTTPError{Title: "Invalid ID", Detail: "An invalid Team ID was provided", Status: http.StatusBadRequest, Cause: err})
		return
	}
	var params renameParams
	err = ctx.ShouldBindJSON(&params)
	if err != nil {
		ctx.Error(&rest.HTTPError{Title: "Invalid Parameters", Detail: "No or invalid parameters where provided to rename the team", Status: http.StatusBadRequest, Cause: err})
		return
	}
	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	err = teamService.RenameTeam(context, uint(id), params.Name)
	cancel()
	if err != nil {
		ctx.Error(&rest.HTTPError{Title: "Unexpected Error", Detail: "An unexpected Error occurde while renaming the team", Status: http.StatusInternalServerError, Cause: err})
		return
	}
	ctx.JSON(http.StatusOK, rest.Response{Data: nil})
}

type recolorParams struct {
	Hue *models.Hue `json:"hue" binding:"required"`
}

func recolorTeam(ctx *gin.Context, teamService *team.TeamService) {
	serializedId := ctx.Param("id")
	id, err := strconv.ParseUint(serializedId, 10, 0)
	if err != nil {
		ctx.Error(&rest.HTTPError{Title: "Invalid ID", Detail: "An invalid Team ID was provided", Status: http.StatusBadRequest, Cause: err})
		return
	}
	var params recolorParams
	err = ctx.ShouldBindJSON(&params)
	if err != nil {
		ctx.Error(&rest.HTTPError{Title: "Invalid Parameters", Detail: "No or invalid parameters where provided to recolor the team", Status: http.StatusBadRequest, Cause: err})
		return
	}
	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	err = teamService.RecolorTeam(context, uint(id), *params.Hue)
	cancel()
	if err != nil {
		ctx.Error(&rest.HTTPError{Title: "Unexpected Error", Detail: "An unexpected Error occurde while recoloring the team", Status: http.StatusInternalServerError, Cause: err})
		return
	}
	ctx.JSON(http.StatusOK, rest.Response{Data: nil})
}

type setOwnerParams struct {
	OwnerID uint `json:"owner_id" binding:"required"`
}

func setOwner(ctx *gin.Context, teamService *team.TeamService) {
	serializedId := ctx.Param("id")
	id, err := strconv.ParseUint(serializedId, 10, 0)
	if err != nil {
		ctx.Error(&rest.HTTPError{Title: "Invalid ID", Detail: "An invalid Team ID was provided", Status: http.StatusBadRequest, Cause: err})
		return
	}
	var params setOwnerParams
	err = ctx.ShouldBindJSON(&params)
	if err != nil {
		ctx.Error(&rest.HTTPError{Title: "Invalid Parameters", Detail: "No or invalid parameters where provided to change the owner of the team", Status: http.StatusBadRequest, Cause: err})
		return
	}
	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	err = teamService.SetOwner(context, uint(id), params.OwnerID)
	cancel()
	if err != nil {
		ctx.Error(&rest.HTTPError{Title: "Unexpected Error", Detail: "An unexpected Error occurde while changing the owner of the team", Status: http.StatusInternalServerError, Cause: err})
		return
	}
	ctx.JSON(http.StatusOK, rest.Response{Data: nil})
}

func getTeams(ctx *gin.Context, teamService *team.TeamService) {
	teams, err := teamService.GetTeams(ctx.Request.Context())
	if err != nil {
		ctx.Error(&rest.HTTPError{Title: "Unexpected Error", Detail: "An unexpected Error occurde while requesting the Teams", Status: http.StatusInternalServerError, Cause: err})
		return
	}
	ctx.JSON(http.StatusOK, rest.Response{Data: teams})
}

func getTeam(ctx *gin.Context, teamService *team.TeamService) {
	serializedId := ctx.Param("id")
	id, err := strconv.ParseUint(serializedId, 10, 0)
	if err != nil {
		ctx.Error(&rest.HTTPError{Title: "Invalid ID", Detail: "An invalid Team ID was provided", Status: http.StatusBadRequest, Cause: err})
		return
	}
	team, err := teamService.GetTeam(ctx.Request.Context(), uint(id))
	if err != nil {
		ctx.Error(&rest.HTTPError{Title: "Unexpected Error", Detail: "An unexpected Error occurde while requesting the Team", Status: http.StatusInternalServerError, Cause: err})
		return
	}
	ctx.JSON(http.StatusOK, rest.Response{Data: team})
}

type inviteParams struct {
	InvitedID uint `json:"invited_id" binding:"required"`
	InviterID uint `json:"inviter_id" binding:"required"`
}

func inviteEntity(ctx *gin.Context, memberService *member.MemberService) {
	serializedId := ctx.Param("id")
	id, err := strconv.ParseUint(serializedId, 10, 0)
	if err != nil {
		ctx.Error(&rest.HTTPError{Title: "Invalid ID", Detail: "An invalid Team ID was provided", Status: http.StatusBadRequest, Cause: err})
		return
	}
	var params inviteParams
	err = ctx.ShouldBindJSON(&params)
	if err != nil {
		ctx.Error(&rest.HTTPError{Title: "Invalid Parameters", Detail: "No or invalid parameters where provided to change invite an entity", Status: http.StatusBadRequest, Cause: err})
		return
	}
	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	invite, err := memberService.CreateInvite(context, params.InvitedID, params.InviterID, uint(id))
	cancel()
	if err != nil {
		ctx.Error(&rest.HTTPError{Title: "Unexpected Error", Detail: "An unexpected Error occurde while creating the invite", Status: http.StatusInternalServerError, Cause: err})
		return
	}
	ctx.JSON(http.StatusOK, rest.Response{Data: invite})
}

func getMembers(ctx *gin.Context, memberService *member.MemberService) {
	serializedId := ctx.Param("id")
	id, err := strconv.ParseUint(serializedId, 10, 0)
	if err != nil {
		ctx.Error(&rest.HTTPError{Title: "Invalid ID", Detail: "An invalid Team ID was provided", Status: http.StatusBadRequest, Cause: err})
		return
	}
	members, err := memberService.GetMembersByTeamId(ctx.Request.Context(), uint(id))
	if err != nil {
		ctx.Error(&rest.HTTPError{Title: "Unexpected Error", Detail: "An unexpected Error occured while requesting the Team Members", Status: http.StatusInternalServerError, Cause: err})
		return
	}
	ctx.JSON(http.StatusOK, rest.Response{Data: members})
}

func getInvites(ctx *gin.Context, memberService *member.MemberService) {
	serializedId := ctx.Param("id")
	id, err := strconv.ParseUint(serializedId, 10, 0)
	if err != nil {
		ctx.Error(&rest.HTTPError{Title: "Invalid ID", Detail: "An invalid Team ID was provided", Status: http.StatusBadRequest, Cause: err})
		return
	}
	members, err := memberService.GetMembersByTeamId(ctx.Request.Context(), uint(id))
	if err != nil {
		ctx.Error(&rest.HTTPError{Title: "Unexpected Error", Detail: "An unexpected Error occured while requesting the Team Members", Status: http.StatusInternalServerError, Cause: err})
		return
	}
	ctx.JSON(http.StatusOK, rest.Response{Data: members})
}
