package handler

import (
	"github.com/annguyen17-tiki/grab/internal/dto"
	"github.com/annguyen17-tiki/grab/internal/model"
	"github.com/annguyen17-tiki/grab/internal/service"
	"github.com/gin-gonic/gin"
)

func saveLocation(svc service.IService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accountID, err := getAccountID(ctx)
		if err != nil {
			renderError(ctx, model.NewErrUnauthorized("unauthorized user"))
			return
		}

		payload := dto.SaveLocation{AccountID: accountID}
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			renderError(ctx, model.NewErrBadRequest(err.Error()))
			return
		}

		if err := globalValidator.Struct(payload); err != nil {
			renderError(ctx, model.NewErrBadRequest(err.Error()))
			return
		}

		if err := svc.SaveLocation(&payload); err != nil {
			renderError(ctx, err)
			return
		}

		renderData(ctx, nil)
	}
}

func nearestLocations(svc service.IService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var queries dto.NearestLocations
		if err := ctx.ShouldBindQuery(&queries); err != nil {
			renderError(ctx, model.NewErrBadRequest(err.Error()))
			return
		}

		if err := globalValidator.Struct(queries); err != nil {
			renderError(ctx, model.NewErrBadRequest(err.Error()))
			return
		}

		locations, err := svc.NearestLocations(&queries)
		if err != nil {
			renderError(ctx, err)
			return
		}

		renderData(ctx, locations)
	}
}
