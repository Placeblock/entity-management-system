package members

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/Placeblock/nostalgicraft-ems/internal/service/member"
	"github.com/Placeblock/nostalgicraft-ems/pkg/models"
	"github.com/Placeblock/nostalgicraft-ems/pkg/rest"
	"github.com/gin-gonic/gin"
)

func Handle(g *gin.RouterGroup, memberService *member.MemberService) {
	g.GET("", func(ctx *gin.Context) {
		getMembers(ctx, memberService)
	})
	g.GET(":id", func(ctx *gin.Context) {
		getMember(ctx, memberService)
	})
	g.DELETE(":id", func(ctx *gin.Context) {
		deleteMember(ctx, memberService)
	})
}

func getMembers(ctx *gin.Context, memberService *member.MemberService) {
	members, err := memberService.GetMembers(ctx.Request.Context())
	if err != nil {
		ctx.Error(&rest.HTTPError{Title: "Unexpected Error", Detail: "An unexpected Error occured while requesting the Team Members", Status: http.StatusInternalServerError, Cause: err})
		return
	}
	ctx.JSON(http.StatusOK, rest.Response{Data: members})
}

func getMember(ctx *gin.Context, memberService *member.MemberService) {
	serializedId := ctx.Param("id")
	id, err := strconv.ParseUint(serializedId, 10, 0)
	if err != nil {
		ctx.Error(&rest.HTTPError{Title: "Invalid ID", Detail: "An invalid Member ID was provided", Status: http.StatusBadRequest, Cause: err})
		return
	}
	members, err := memberService.GetMember(ctx.Request.Context(), uint(id))
	if err != nil {
		ctx.Error(&rest.HTTPError{Title: "Unexpected Error", Detail: "An unexpected Error occured while requesting the Team Members", Status: http.StatusInternalServerError, Cause: err})
		return
	}
	ctx.JSON(http.StatusOK, rest.Response{Data: members})
}

func deleteMember(ctx *gin.Context, memberService *member.MemberService) {
	serializedId := ctx.Param("id")
	id, err := strconv.ParseUint(serializedId, 10, 0)
	if err != nil {
		ctx.Error(&rest.HTTPError{Title: "Invalid ID", Detail: "An invalid Member ID was provided", Status: http.StatusBadRequest, Cause: err})
		return
	}
	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	err = memberService.RemoveMember(context, &models.Member{ID: uint(id)})
	cancel()
	if err != nil {
		ctx.Error(&rest.HTTPError{Title: "Unexpected Error", Detail: "An unexpected Error occured while deleting the Team Member", Status: http.StatusInternalServerError, Cause: err})
		return
	}
	ctx.JSON(http.StatusOK, rest.Response{Data: nil})
}
