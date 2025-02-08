package teams

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/codelix/ems/internal/service/team"
	"github.com/codelix/ems/pkg/models"
	"github.com/codelix/ems/pkg/rest"
	"github.com/gin-gonic/gin"
)

func Handle(g *gin.RouterGroup, teamService *team.TeamService) {
	g.GET("", func(ctx *gin.Context) {
		getTeams(ctx, teamService)
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
	g.POST("", func(ctx *gin.Context) {
		createTeam(ctx, teamService)
	})
}

type createParams struct {
	Name string     `json:"name" binding:"required"`
	Hue  models.Hue `json:"hue" binding:"required"`
}

func createTeam(ctx *gin.Context, teamService *team.TeamService) {
	var params createParams
	err := ctx.ShouldBindJSON(&params)
	if err != nil {
		ctx.Error(&rest.HTTPError{Title: "Invalid Parameters", Detail: "No or invalid parameters where provided to create the team", Status: http.StatusBadRequest, Cause: err})
		return
	}
	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	var team models.Team
	err = teamService.CreateTeam(context, &team)
	cancel()
	if err != nil {
		ctx.Error(&rest.HTTPError{Title: "Unexpected Error", Detail: "An unexpected Error occurde while creating the team", Status: http.StatusInternalServerError, Cause: err})
		return
	}
	ctx.JSON(http.StatusOK, team)
}

type renameParams struct {
	Name string `json:"name" binding:"required"`
}

func renameTeam(ctx *gin.Context, teamService *team.TeamService) {
	serializedId := ctx.Param("id")
	id, err := strconv.ParseUint(serializedId, 10, 0)
	if err != nil {
		ctx.Error(&rest.HTTPError{Title: "Invalid ID", Detail: "No or an invalid Team ID was provided", Status: http.StatusBadRequest, Cause: err})
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
	ctx.JSON(http.StatusOK, nil)
}

type recolorParams struct {
	Hue models.Hue `json:"hue" binding:"required"`
}

func recolorTeam(ctx *gin.Context, teamService *team.TeamService) {
	serializedId := ctx.Param("id")
	id, err := strconv.ParseUint(serializedId, 10, 0)
	if err != nil {
		ctx.Error(&rest.HTTPError{Title: "Invalid ID", Detail: "No or an invalid Team ID was provided", Status: http.StatusBadRequest, Cause: err})
		return
	}
	var params recolorParams
	err = ctx.ShouldBindJSON(&params)
	if err != nil {
		ctx.Error(&rest.HTTPError{Title: "Invalid Parameters", Detail: "No or invalid parameters where provided to recolor the team", Status: http.StatusBadRequest, Cause: err})
		return
	}
	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	err = teamService.RecolorTeam(context, uint(id), params.Hue)
	cancel()
	if err != nil {
		ctx.Error(&rest.HTTPError{Title: "Unexpected Error", Detail: "An unexpected Error occurde while recoloring the team", Status: http.StatusInternalServerError, Cause: err})
		return
	}
	ctx.JSON(http.StatusOK, nil)
}

type setOwnerParams struct {
	OwnerID uint `json:"owner_id" binding:"required"`
}

func setOwner(ctx *gin.Context, teamService *team.TeamService) {
	serializedId := ctx.Param("id")
	id, err := strconv.ParseUint(serializedId, 10, 0)
	if err != nil {
		ctx.Error(&rest.HTTPError{Title: "Invalid ID", Detail: "No or an invalid Team ID was provided", Status: http.StatusBadRequest, Cause: err})
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
	ctx.JSON(http.StatusOK, nil)
}

func getTeams(ctx *gin.Context, teamService *team.TeamService) {
	teams, err := teamService.GetTeams(ctx.Request.Context())
	if err != nil {
		ctx.Error(&rest.HTTPError{Title: "Unexpected Error", Detail: "An unexpected Error occurde while requesting the Teams", Status: http.StatusInternalServerError, Cause: err})
		return
	}
	ctx.JSON(http.StatusOK, teams)
}

func getTeam(ctx *gin.Context, teamService *team.TeamService) {
	serializedId := ctx.Param("id")
	id, err := strconv.ParseUint(serializedId, 10, 0)
	if err != nil {
		ctx.Error(&rest.HTTPError{Title: "Invalid ID", Detail: "No or an invalid Team ID was provided", Status: http.StatusBadRequest, Cause: err})
		return
	}
	team, err := teamService.GetTeam(ctx.Request.Context(), uint(id))
	if err != nil {
		ctx.Error(&rest.HTTPError{Title: "Unexpected Error", Detail: "An unexpected Error occurde while requesting the Team", Status: http.StatusInternalServerError, Cause: err})
		return
	}
	ctx.JSON(http.StatusOK, team)
}
