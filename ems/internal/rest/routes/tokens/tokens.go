package tokens

import (
	"net/http"

	"github.com/Placeblock/nostalgicraft-ems/internal/service/token"
	"github.com/Placeblock/nostalgicraft-ems/pkg/rest"
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
		ctx.Error(&rest.HTTPError{Title: "Invalid Pin", Detail: "No or an invalid Pin was provided", Status: http.StatusBadRequest})
		return
	}
	token, err := tokenService.GetToken(ctx.Request.Context(), pin)
	if err != nil {
		ctx.Error(&rest.HTTPError{Title: "Unexpected Error", Detail: "An unexpected Error occurde while requesting the token", Status: http.StatusInternalServerError, Cause: err})
		return
	}
	ctx.JSON(http.StatusOK, rest.Response{Data: token})
}

type createTokenParams struct {
	EntityId uint `json:"entityId" binding:"required"`
}

func createToken(ctx *gin.Context, tokenService *token.TokenService) {
	var params createTokenParams
	err := ctx.ShouldBindJSON(&params)
	if err != nil {
		ctx.Error(&rest.HTTPError{Title: "Invalid Token Parameters", Detail: "No or invalid parameters where provided to create the token", Status: http.StatusBadRequest})
		return
	}
	token, err := tokenService.CreateToken(ctx.Request.Context(), params.EntityId)
	if err != nil {
		ctx.Error(&rest.HTTPError{Title: "Unexpected Error", Detail: "An unexpected Error occurde while creating the token", Status: http.StatusInternalServerError, Cause: err})
		return
	}
	ctx.JSON(http.StatusOK, rest.Response{Data: token})
}
