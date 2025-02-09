package invites

import (
	"github.com/codelix/ems/internal/service/member"
	"github.com/gin-gonic/gin"
)

func Handle(g *gin.RouterGroup, memberService *member.MemberService) {
	g.GET("", func(ctx *gin.Context) {
		getInvites(ctx, memberService)
	})
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
