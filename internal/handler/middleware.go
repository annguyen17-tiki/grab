package handler

import (
	"errors"
	"strings"

	"github.com/annguyen17-tiki/grab/internal/model"
	"github.com/annguyen17-tiki/grab/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

const (
	ctxKeyAccountID = "account_id"
	ctxKeyRole      = "role"
)

func authenticate(ctx *gin.Context) {
	bearerHeader := ctx.Request.Header.Get("Authorization")
	if bearerHeader == "" {
		renderError(ctx, model.NewErrUnauthorized("missing token"))
		return
	}

	parts := strings.Split(bearerHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		renderError(ctx, model.NewErrUnauthorized("invalid token"))
		return
	}

	var claim jwt.StandardClaims
	if _, err := jwt.ParseWithClaims(parts[1], &claim, func(token *jwt.Token) (interface{}, error) {
		return []byte(getConfig().JWTSecret), nil
	}); err != nil {
		renderError(ctx, model.NewErrUnauthorized(err.Error()))
		return
	}

	ctx.Set(ctxKeyAccountID, claim.Subject)
	ctx.Next()
}

func requireRole(svc service.IService, roles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accountID, err := getAccountID(ctx)
		if err != nil {
			renderError(ctx, model.NewErrUnauthorized("unauthorized user"))
			return
		}

		account, err := svc.GetAccount(accountID)
		if err != nil {
			renderError(ctx, err)
			return
		}

		hasRole := false
		for _, r := range roles {
			if account.Role == r {
				hasRole = true
				break
			}
		}

		if !hasRole {
			renderError(ctx, model.NewErrForbidden("require roles: %s", strings.Join(roles, ", ")))
			return
		}

		ctx.Next()
	}
}

func getAccountID(ctx *gin.Context) (string, error) {
	iAccountID, found := ctx.Get(ctxKeyAccountID)
	if !found {
		return "", errors.New("account id not found")
	}

	accountID, found := iAccountID.(string)
	if !found {
		return "", errors.New("account id not found")
	}

	return accountID, nil
}
