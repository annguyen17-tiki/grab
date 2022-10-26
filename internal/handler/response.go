package handler

import (
	"net/http"

	"github.com/annguyen17-tiki/grab/internal/model"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

func renderData(ctx *gin.Context, data interface{}) {
	ctx.AbortWithStatusJSON(http.StatusOK, Response{Data: data})
}

func renderError(ctx *gin.Context, err error) {
	switch err.(type) {
	case *model.ErrBadRequest:
		ctx.AbortWithStatusJSON(http.StatusBadRequest, Response{Error: err.Error()})
		return
	case *model.ErrUnauthorized:
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, Response{Error: err.Error()})
		return
	case *model.ErrForbidden:
		ctx.AbortWithStatusJSON(http.StatusForbidden, Response{Error: err.Error()})
		return
	case *model.ErrNotFound:
		ctx.AbortWithStatusJSON(http.StatusNotFound, Response{Error: err.Error()})
		return
	default:
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}
}
