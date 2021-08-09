package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/error"
	"github.com/ictsc/ictsc-rikka/pkg/repository"
)

func Auth(userRepo repository.UserRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)

		idString, ok := session.Get("id").(string)
		if !ok || idString == "" {
			ctx.Error(error.NewUnauthorizedError("unauthorized"))
			ctx.Abort()
			return
		}

		id, err := uuid.Parse(idString)
		if err != nil {
			ctx.Error(error.NewInternalServerError(err))
			ctx.Abort()
			return
		}

		user, err := userRepo.FindByID(id, true)
		if err != nil {
			ctx.Error(error.NewInternalServerError(err))
			ctx.Abort()
			return
		}

		ctx.Set("is_full_access", user.UserGroup.IsFullAccess)
		ctx.Set("user", user)
		ctx.Set("group", user.UserGroup)

		ctx.Next()
	}
}

func AuthIsFullAccess(userRepo repository.UserRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)

		idString, ok := session.Get("id").(string)
		if !ok || idString == "" {
			ctx.Error(error.NewUnauthorizedError("unauthorized"))
			ctx.Abort()
			return
		}

		id, err := uuid.Parse(idString)
		if err != nil {
			ctx.Error(error.NewInternalServerError(err))
			ctx.Abort()
			return
		}

		user, err := userRepo.FindByID(id, true)
		if err != nil {
			ctx.Error(error.NewInternalServerError(err))
			ctx.Abort()
			return
		}

		if user.UserGroup == nil {
			panic("user.UserGroup must not be nil")
		}

		if !user.UserGroup.IsFullAccess {
			ctx.Error(error.NewForbiddenError("you don't have enough permission"))
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
