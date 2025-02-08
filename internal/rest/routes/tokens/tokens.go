package tokens

import (
	"net/http"

	"github.com/codelix/ems/internal/service/token"
	"github.com/gin-gonic/gin"
)

func Handle(g *gin.RouterGroup, service *token.TokenService) {
	g.GET("", func(ctx *gin.Context) {
		getToken(ctx, service)
	})
	g.POST("", func(ctx *gin.Context) {
		createToken(ctx, service)
	})
}

func getToken(ctx *gin.Context, tokenService *token.TokenService) {
	pin := ctx.Query("pin")
	if pin == "" {
		ctx.JSON(http.StatusBadRequest, "Invalid Pin")
		return
	}
	token, err := tokenService.GetToken(ctx.Request.Context(), pin)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, token)
}

func createToken(ctx *gin.Context, tokenService *token.TokenService) {
	var params createTokenParams
	err := ctx.ShouldBindJSON(&params)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	token, err := tokenService.CreateToken(ctx.Request.Context(), params.EntityId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, token)
}
