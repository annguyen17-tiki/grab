package handler

import (
	"github.com/annguyen17-tiki/grab/internal/dto"
	"github.com/annguyen17-tiki/grab/internal/model"
	"github.com/annguyen17-tiki/grab/internal/service"
	"github.com/gin-gonic/gin"
)

func createAccount(svc service.IService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload dto.CreateAccount
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			renderError(ctx, model.NewErrBadRequest(err.Error()))
			return
		}

		if err := globalValidator.Struct(payload); err != nil {
			renderError(ctx, model.NewErrBadRequest(err.Error()))
			return
		}

		if err := svc.CreateAccount(&payload); err != nil {
			renderError(ctx, err)
			return
		}

		renderData(ctx, nil)
	}
}

func updateAccount(svc service.IService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accountID, err := getAccountID(ctx)
		if err != nil {
			renderError(ctx, model.NewErrForbidden(err.Error()))
			return
		}

		payload := dto.UpdateAccount{ID: accountID}
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			renderError(ctx, model.NewErrBadRequest(err.Error()))
			return
		}

		if err := globalValidator.Struct(payload); err != nil {
			renderError(ctx, model.NewErrBadRequest(err.Error()))
			return
		}

		if err := svc.UpdateAccount(&payload); err != nil {
			renderError(ctx, err)
			return
		}

		renderData(ctx, nil)
	}
}

func login(svc service.IService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload dto.LoginInput
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			renderError(ctx, model.NewErrBadRequest(err.Error()))
			return
		}

		if err := globalValidator.Struct(payload); err != nil {
			renderError(ctx, model.NewErrBadRequest(err.Error()))
			return
		}

		token, err := svc.Login(&payload)
		if err != nil {
			renderError(ctx, err)
			return
		}

		renderData(ctx, token)
	}
}

func getOwnAccount(svc service.IService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accountID, err := getAccountID(ctx)
		if err != nil {
			renderError(ctx, model.NewErrForbidden(err.Error()))
			return
		}

		account, err := svc.GetAccount(accountID)
		if err != nil {
			renderError(ctx, err)
			return
		}

		renderData(ctx, account)
	}
}

func getAccount(svc service.IService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		phone := ctx.Query("phone")
		if err := globalValidator.Var(phone, "required,phone"); err != nil {
			renderError(ctx, model.NewErrBadRequest(err.Error()))
			return
		}

		account, err := svc.GetAccountByPhone(phone)
		if err != nil {
			renderError(ctx, err)
			return
		}

		renderData(ctx, account)
	}
}
