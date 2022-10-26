package handler

import (
	"github.com/annguyen17-tiki/grab/internal/dto"
	"github.com/annguyen17-tiki/grab/internal/model"
	"github.com/annguyen17-tiki/grab/internal/service"
	"github.com/gin-gonic/gin"
)

func createBooking(svc service.IService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accountID, err := getAccountID(ctx)
		if err != nil {
			renderError(ctx, model.NewErrUnauthorized("unauthorized user"))
			return
		}

		payload := dto.CreateBooking{UserID: accountID}
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			renderError(ctx, model.NewErrBadRequest(err.Error()))
			return
		}

		if err := globalValidator.Struct(payload); err != nil {
			renderError(ctx, model.NewErrBadRequest(err.Error()))
			return
		}

		booking, err := svc.CreateBooking(&payload)
		if err != nil {
			renderError(ctx, err)
			return
		}

		renderData(ctx, booking)
	}
}

func acceptBooking(svc service.IService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accountID, err := getAccountID(ctx)
		if err != nil {
			renderError(ctx, model.NewErrUnauthorized("unauthorized user"))
			return
		}

		bookingID := ctx.Params.ByName("booking_id")
		if err := globalValidator.Var(bookingID, "uuid"); err != nil {
			renderError(ctx, model.NewErrBadRequest(err.Error()))
			return
		}

		if err := svc.AcceptBooking(bookingID, accountID); err != nil {
			renderError(ctx, err)
			return
		}

		renderData(ctx, nil)
	}
}

func rejectBooking(svc service.IService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accountID, err := getAccountID(ctx)
		if err != nil {
			renderError(ctx, model.NewErrUnauthorized("unauthorized user"))
			return
		}

		bookingID := ctx.Params.ByName("booking_id")
		if err := globalValidator.Var(bookingID, "uuid"); err != nil {
			renderError(ctx, model.NewErrBadRequest(err.Error()))
			return
		}

		if err := svc.RejectBooking(bookingID, accountID); err != nil {
			renderError(ctx, err)
			return
		}

		renderData(ctx, nil)
	}
}

func doneBooking(svc service.IService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accountID, err := getAccountID(ctx)
		if err != nil {
			renderError(ctx, model.NewErrUnauthorized("unauthorized user"))
			return
		}

		bookingID := ctx.Params.ByName("booking_id")
		if err := globalValidator.Var(bookingID, "uuid"); err != nil {
			renderError(ctx, model.NewErrBadRequest(err.Error()))
			return
		}

		if err := svc.DoneBooking(bookingID, accountID); err != nil {
			renderError(ctx, err)
			return
		}

		renderData(ctx, nil)
	}
}

func getBooking(svc service.IService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bookingID := ctx.Params.ByName("booking_id")
		if err := globalValidator.Var(bookingID, "uuid"); err != nil {
			renderError(ctx, model.NewErrBadRequest(err.Error()))
			return
		}

		booking, err := svc.GetBooking(bookingID)
		if err != nil {
			renderError(ctx, err)
			return
		}

		renderData(ctx, booking)
	}
}

func searchBookings(svc service.IService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		queries := dto.SearchBooking{Paging: &defaultPagination}
		if err := ctx.ShouldBindQuery(&queries); err != nil {
			renderError(ctx, model.NewErrBadRequest(err.Error()))
			return
		}

		accountID, err := getAccountID(ctx)
		if err != nil {
			renderError(ctx, model.NewErrUnauthorized("unauthorized user"))
			return
		}

		queries.AccountID = &accountID

		if err := globalValidator.Struct(queries); err != nil {
			renderError(ctx, model.NewErrBadRequest(err.Error()))
			return
		}

		bookings, err := svc.SearchBooking(&queries)
		if err != nil {
			renderError(ctx, err)
			return
		}

		renderData(ctx, bookings)
	}
}

func searchBookingsForAdmin(svc service.IService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		queries := dto.SearchBooking{Paging: &defaultPagination}
		if err := ctx.ShouldBindQuery(&queries); err != nil {
			renderError(ctx, model.NewErrBadRequest(err.Error()))
			return
		}

		if err := globalValidator.Struct(queries); err != nil {
			renderError(ctx, model.NewErrBadRequest(err.Error()))
			return
		}

		bookings, err := svc.SearchBooking(&queries)
		if err != nil {
			renderError(ctx, err)
			return
		}

		renderData(ctx, bookings)
	}
}
