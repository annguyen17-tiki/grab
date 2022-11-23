package handler

import (
	"github.com/annguyen17-tiki/grab/internal/dto"
	"github.com/annguyen17-tiki/grab/internal/model"
	"github.com/annguyen17-tiki/grab/internal/service"
	"github.com/gin-gonic/gin"
)

func saveToken(svc service.IService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accountID, err := getAccountID(ctx)
		if err != nil {
			renderError(ctx, model.NewErrUnauthorized("unauthorized user"))
			return
		}

		var payload dto.SaveToken
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			renderError(ctx, model.NewErrBadRequest(err.Error()))
			return
		}

		if err := globalValidator.Struct(payload); err != nil {
			renderError(ctx, model.NewErrBadRequest(err.Error()))
			return
		}

		if err := svc.SaveToken(accountID, payload.Token); err != nil {
			renderError(ctx, err)
			return
		}

		renderData(ctx, nil)
	}
}
