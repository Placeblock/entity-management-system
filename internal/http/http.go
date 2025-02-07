package http

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/codelix/ems/internal/service/entity"
	"github.com/codelix/ems/internal/service/token"
	"github.com/codelix/ems/pkg/models"
	"github.com/gin-gonic/gin"
)

type updateEntityParams struct {
	NewName string `json:"new_name" binding:"required"`
}

type createEntityParams struct {
	Name string `json:"name" binding:"required"`
}

type createTokenParams struct {
	EntityId int64 `json:"entityId" binding:"required"`
}

type HttpServer struct {
	entityService entity.EntityService
	tokenService  token.TokenService
}

func NewHttpServer(entityService entity.EntityService, tokenService token.TokenService) *HttpServer {
	return &HttpServer{entityService, tokenService}
}

func (server *HttpServer) Serve() {
	r := gin.Default()
	entitiesGroup := r.Group("entities")

	// GET ENTITY
	entitiesGroup.GET(":id", func(ctx *gin.Context) {
		serializedId := ctx.Param("id")
		id, err := strconv.ParseInt(serializedId, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, "Invalid id")
			return
		}
		entity, err := server.entityService.GetEntity(ctx.Request.Context(), id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, entity)
	})

	// UPDATE ENTITY
	entitiesGroup.PUT(":id", func(ctx *gin.Context) {
		serializedId := ctx.Param("id")
		id, err := strconv.ParseInt(serializedId, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, "Invalid id")
			return
		}
		var params updateEntityParams
		err = ctx.ShouldBindJSON(&params)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}
		context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err = server.entityService.RenameEntity(context, id, params.NewName)
		cancel()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, nil)
	})

	// GET ALL ENTITIES
	entitiesGroup.GET("", func(ctx *gin.Context) {
		entities, err := server.entityService.GetEntities(ctx.Request.Context())
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, entities)
	})

	// CREATE ENTITY
	entitiesGroup.POST("", func(ctx *gin.Context) {
		var params createEntityParams
		err := ctx.ShouldBindJSON(&params)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}
		context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		entity := models.Entity{Name: params.Name}
		err = server.entityService.CreateEntity(context, &entity)
		cancel()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, entity)
	})

	tokenGroup := r.Group("token")

	tokenGroup.GET("", func(ctx *gin.Context) {
		pin := ctx.Query("pin")
		if pin == "" {
			ctx.JSON(http.StatusBadRequest, "Invalid Pin")
			return
		}
		token, err := server.tokenService.GetToken(ctx.Request.Context(), pin)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, token)
	})

	tokenGroup.POST("", func(ctx *gin.Context) {
		var params createTokenParams
		err := ctx.ShouldBindJSON(&params)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}
		token, err := server.tokenService.CreateToken(ctx.Request.Context(), params.EntityId)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, token)
	})

	r.Run()
}
