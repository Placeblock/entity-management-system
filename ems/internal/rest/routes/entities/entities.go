package entities

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/Placeblock/nostalgicraft-ems/internal/service/entity"
	"github.com/Placeblock/nostalgicraft-ems/internal/service/member"
	"github.com/Placeblock/nostalgicraft-ems/pkg/models"
	"github.com/Placeblock/nostalgicraft-ems/pkg/rest"
	"github.com/gin-gonic/gin"
)

func Handle(g *gin.RouterGroup, entityService *entity.EntityService, memberService *member.MemberService) {
	g.POST("", func(ctx *gin.Context) {
		createEntity(ctx, entityService)
	})
	g.GET("", func(ctx *gin.Context) {
		getEntities(ctx, entityService)
	})
	g.GET(":id", func(ctx *gin.Context) {
		getEntity(ctx, entityService)
	})
	g.GET(":id/member", func(ctx *gin.Context) {
		getMember(ctx, memberService)
	})
	g.PUT(":id", func(ctx *gin.Context) {
		renameEntity(ctx, entityService)
	})
	g.DELETE(":id/team", func(ctx *gin.Context) {
		leaveTeam(ctx, memberService)
	})
	g.GET(":id/invites", func(ctx *gin.Context) {
		getInvitesByInvited(ctx, memberService)
	})
}

func getEntity(ctx *gin.Context, entityService *entity.EntityService) {
	serializedId := ctx.Param("id")
	id, err := strconv.ParseUint(serializedId, 10, 0)
	if err != nil {
		ctx.Error(&rest.HTTPError{Title: "Invalid ID", Detail: "No or an invalid Entity ID was provided", Status: http.StatusBadRequest, Cause: err})
		return
	}
	entity, err := entityService.GetEntity(ctx.Request.Context(), uint(id))
	if err != nil {
		ctx.Error(&rest.HTTPError{Title: "Unexpected Error", Detail: "An unexpected Error occurde while requesting the Entity", Status: http.StatusInternalServerError, Cause: err})
		return
	}
	ctx.JSON(http.StatusOK, rest.Response{Data: entity})
}

func getMember(ctx *gin.Context, memberService *member.MemberService) {
	serializedId := ctx.Param("id")
	id, err := strconv.ParseUint(serializedId, 10, 0)
	if err != nil {
		ctx.Error(&rest.HTTPError{Title: "Invalid ID", Detail: "No or an invalid Entity ID was provided", Status: http.StatusBadRequest, Cause: err})
		return
	}
	member, err := memberService.GetMemberByEntityId(ctx.Request.Context(), uint(id))
	if err != nil {
		ctx.Error(&rest.HTTPError{Title: "Unexpected Error", Detail: "An unexpected Error occurde while requesting the Member", Status: http.StatusInternalServerError, Cause: err})
		return
	}
	ctx.JSON(http.StatusOK, rest.Response{Data: member})
}

func getEntities(ctx *gin.Context, entityService *entity.EntityService) {
	entities, err := entityService.GetEntities(ctx.Request.Context())
	if err != nil {
		ctx.Error(&rest.HTTPError{Title: "Unexpected Error", Detail: "An unexpected Error occurde while requesting the Entities", Status: http.StatusInternalServerError, Cause: err})
		return
	}
	ctx.JSON(http.StatusOK, rest.Response{Data: entities})
}

type createEntityParams struct {
	Name string `json:"name" binding:"required"`
}

func createEntity(ctx *gin.Context, entityService *entity.EntityService) {
	var params createEntityParams
	err := ctx.ShouldBindJSON(&params)
	if err != nil {
		ctx.Error(&rest.HTTPError{Title: "Invalid Entity Parameters", Detail: "No or invalid parameters where provided to create the Entity", Status: http.StatusBadRequest, Cause: err})
		return
	}
	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	entity := models.Entity{Name: params.Name}
	err = entityService.CreateEntity(context, &entity)
	cancel()
	if err != nil {
		ctx.Error(&rest.HTTPError{Title: "Unexpected Error", Detail: "An unexpected Error occurde while creating the Entity", Status: http.StatusInternalServerError, Cause: err})
		return
	}
	ctx.JSON(http.StatusOK, rest.Response{Data: entity})
}

type updateEntityParams struct {
	NewName string `json:"new_name" binding:"required"`
}

func renameEntity(ctx *gin.Context, entityService *entity.EntityService) {
	serializedId := ctx.Param("id")
	id, err := strconv.ParseUint(serializedId, 10, 0)
	if err != nil {
		ctx.Error(&rest.HTTPError{Title: "Invalid ID", Detail: "No or an invalid Entity ID was provided", Status: http.StatusBadRequest, Cause: err})
		return
	}
	var params updateEntityParams
	err = ctx.ShouldBindJSON(&params)
	if err != nil {
		ctx.Error(&rest.HTTPError{Title: "Invalid Entity Parameters", Detail: "No or invalid parameters where provided to rename the Entity", Status: http.StatusBadRequest, Cause: err})
		return
	}
	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	err = entityService.RenameEntity(context, uint(id), params.NewName)
	cancel()
	if err != nil {
		ctx.Error(&rest.HTTPError{Title: "Unexpected Error", Detail: "An unexpected Error occurde while renaming the Entity", Status: http.StatusInternalServerError, Cause: err})
		return
	}
	ctx.JSON(http.StatusOK, rest.Response{Data: nil})
}

func leaveTeam(ctx *gin.Context, memberService *member.MemberService) {
	serializedId := ctx.Param("id")
	id, err := strconv.ParseUint(serializedId, 10, 0)
	if err != nil {
		ctx.Error(&rest.HTTPError{Title: "Invalid ID", Detail: "No or an invalid Entity ID was provided", Status: http.StatusBadRequest, Cause: err})
		return
	}
	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	err = memberService.RemoveMember(context, &models.Member{EntityID: uint(id)})
	cancel()
	if err != nil {
		ctx.Error(&rest.HTTPError{Title: "Unexpected Error", Detail: "An unexpected Error occurde while leaving the Team", Status: http.StatusInternalServerError, Cause: err})
		return
	}
	ctx.JSON(http.StatusOK, rest.Response{Data: nil})
}

func getInvitesByInvited(ctx *gin.Context, memberService *member.MemberService) {
	serializedId := ctx.Param("id")
	id, err := strconv.ParseUint(serializedId, 10, 0)
	if err != nil {
		ctx.Error(&rest.HTTPError{Title: "Invalid ID", Detail: "No or an invalid Entity ID was provided", Status: http.StatusBadRequest, Cause: err})
		return
	}
	invites, err := memberService.GetMemberInvitesByInvitedId(ctx, uint(id))
	if err != nil {
		ctx.Error(&rest.HTTPError{Title: "Unexpected Error", Detail: "An unexpected Error occurde while requesting invites", Status: http.StatusInternalServerError, Cause: err})
		return
	}
	ctx.JSON(http.StatusOK, rest.Response{Data: invites})
}
