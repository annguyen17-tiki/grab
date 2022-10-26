package handler

import (
	"github.com/annguyen17-tiki/grab/internal/dto"
	"github.com/annguyen17-tiki/grab/internal/model"
	"github.com/annguyen17-tiki/grab/internal/service"
	"github.com/gin-gonic/gin"
)

func searchNotifications(svc service.IService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accountID, err := getAccountID(ctx)
		if err != nil {
			renderError(ctx, model.NewErrUnauthorized("unauthorized user"))
			return
		}

		queries := dto.SearchNotification{AccountID: &accountID, Paging: &defaultPagination}
		if err := ctx.ShouldBindQuery(&queries); err != nil {
			renderError(ctx, model.NewErrBadRequest(err.Error()))
			return
		}

		if err := globalValidator.Struct(queries); err != nil {
			renderError(ctx, model.NewErrBadRequest(err.Error()))
			return
		}

		notifications, err := svc.SearchNotifications(&queries)
		if err != nil {
			renderError(ctx, err)
			return
		}

		renderData(ctx, notifications)
	}
}

func seenNotification(svc service.IService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accountID, err := getAccountID(ctx)
		if err != nil {
			renderError(ctx, model.NewErrUnauthorized("unauthorized user"))
			return
		}

		notificationID := ctx.Params.ByName("notification_id")
		if err := globalValidator.Var(notificationID, "uuid"); err != nil {
			renderError(ctx, model.NewErrBadRequest(err.Error()))
			return
		}

		if err := svc.SeenNotification(accountID, notificationID); err != nil {
			renderError(ctx, err)
			return
		}

		renderData(ctx, nil)
	}
}
