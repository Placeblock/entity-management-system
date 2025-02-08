package entities

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/codelix/ems/internal/service/entity"
	"github.com/codelix/ems/pkg/models"
	"github.com/gin-gonic/gin"
)

func Handle(g *gin.RouterGroup, entityService *entity.EntityService) {
	g.GET(":id", func(ctx *gin.Context) {
		getEntity(ctx, entityService)
	})
	g.PUT(":id", func(ctx *gin.Context) {
		renameEntity(ctx, entityService)
	})
	g.GET("", func(ctx *gin.Context) {
		getEntities(ctx, entityService)
	})
	g.POST("", func(ctx *gin.Context) {
		createEntity(ctx, entityService)
	})
}

func getEntity(ctx *gin.Context, entityService *entity.EntityService) {
	serializedId := ctx.Param("id")
	id, err := strconv.ParseUint(serializedId, 10, 0)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "Invalid id")
		return
	}
	entity, err := entityService.GetEntity(ctx.Request.Context(), uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, entity)
}

func getEntities(ctx *gin.Context, entityService *entity.EntityService) {
	entities, err := entityService.GetEntities(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, entities)
}

func createEntity(ctx *gin.Context, entityService *entity.EntityService) {
	var params createEntityParams
	err := ctx.ShouldBindJSON(&params)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	entity := models.Entity{Name: params.Name}
	err = entityService.CreateEntity(context, &entity)
	cancel()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, entity)
}

func renameEntity(ctx *gin.Context, entityService *entity.EntityService) {
	serializedId := ctx.Param("id")
	id, err := strconv.ParseUint(serializedId, 10, 0)
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
	err = entityService.RenameEntity(context, uint(id), params.NewName)
	cancel()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, nil)
}
