package invites

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/Placeblock/nostalgicraft-ems/internal/service/member"
	"github.com/Placeblock/nostalgicraft-ems/pkg/rest"
	"github.com/gin-gonic/gin"
)

func Handle(g *gin.RouterGroup, memberService *member.MemberService) {
	g.GET(":id", func(ctx *gin.Context) {
		getInvite(ctx, memberService)
	})
	g.POST(":id", func(ctx *gin.Context) {
		acceptInvite(ctx, memberService)
	})
	g.DELETE(":id", func(ctx *gin.Context) {
		declineInvite(ctx, memberService)
	})
	g.POST("", func(ctx *gin.Context) {
		createInvite(ctx, memberService)
	})
}

type createInviteParams struct {
	InvitedID uint `json:"invited_id" binding:"required"`
	InviterID uint `json:"inviter_id" binding:"required"`
	TeamID    uint `json:"team_id" binding:"required"`
}

func createInvite(ctx *gin.Context, memberService *member.MemberService) {
	var params createInviteParams
	err := ctx.ShouldBindJSON(&params)
	if err != nil {
		ctx.Error(&rest.HTTPError{Title: "Invalid Invite Parameters", Detail: "No or invalid parameters where provided to create the Invite", Status: http.StatusBadRequest, Cause: err})
		return
	}
	if params.InvitedID == params.InviterID {
		ctx.Error(&rest.HTTPError{Title: "Invalid Invite Parameters", Detail: "The inviter and the invited cannot be the same entity", Status: http.StatusBadRequest, Cause: err})
		return
	}
	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	invite, err := memberService.CreateInvite(context, params.InvitedID, params.InviterID, params.TeamID)
	cancel()
	if err != nil {
		ctx.Error(&rest.HTTPError{Title: "Unexpected Error", Detail: "An unexpected Error occurde while creating the Invite", Status: http.StatusInternalServerError, Cause: err})
		return
	}
	ctx.JSON(http.StatusOK, rest.Response{Data: invite})
}

func declineInvite(ctx *gin.Context, memberService *member.MemberService) {
	serializedId := ctx.Param("id")
	id, err := strconv.ParseUint(serializedId, 10, 0)
	if err != nil {
		ctx.Error(&rest.HTTPError{Title: "Invalid ID", Detail: "No or an invalid Invite ID was provided", Status: http.StatusBadRequest, Cause: err})
		return
	}
	err = memberService.DeclineInvite(ctx.Request.Context(), uint(id))
	if err != nil {
		ctx.Error(&rest.HTTPError{Title: "Unexpected Error", Detail: "An unexpected Error occurde while declining the Invite", Status: http.StatusInternalServerError, Cause: err})
		return
	}
	ctx.JSON(http.StatusOK, rest.Response{Data: nil})
}

func acceptInvite(ctx *gin.Context, memberService *member.MemberService) {
	serializedId := ctx.Param("id")
	id, err := strconv.ParseUint(serializedId, 10, 0)
	if err != nil {
		ctx.Error(&rest.HTTPError{Title: "Invalid ID", Detail: "No or an invalid Invite ID was provided", Status: http.StatusBadRequest, Cause: err})
		return
	}
	member, err := memberService.AcceptInvite(ctx.Request.Context(), uint(id))
	if err != nil {
		ctx.Error(&rest.HTTPError{Title: "Unexpected Error", Detail: "An unexpected Error occurde while accepting the Invite", Status: http.StatusInternalServerError, Cause: err})
		return
	}
	ctx.JSON(http.StatusOK, rest.Response{Data: member})
}

func getInvite(ctx *gin.Context, memberService *member.MemberService) {
	serializedId := ctx.Param("id")
	id, err := strconv.ParseUint(serializedId, 10, 0)
	if err != nil {
		ctx.Error(&rest.HTTPError{Title: "Invalid ID", Detail: "No or an invalid Entity ID was provided", Status: http.StatusBadRequest, Cause: err})
		return
	}
	invite, err := memberService.GetMemberInvite(ctx.Request.Context(), uint(id))
	if err != nil {
		ctx.Error(&rest.HTTPError{Title: "Unexpected Error", Detail: "An unexpected Error occurde while requesting the Invite", Status: http.StatusInternalServerError, Cause: err})
		return
	}
	ctx.JSON(http.StatusOK, rest.Response{Data: invite})
}
