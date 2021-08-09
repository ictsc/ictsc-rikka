package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/repository"
)

func Auth(userRepo repository.UserRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)

		idString, ok := session.Get("id").(string)
		if !ok || idString == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "unauthorized.",
			})
			return
		}

		id, err := uuid.Parse(idString)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "internal server error",
			})
			return
		}

		user, err := userRepo.FindByID(id, true)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "internal server error",
			})
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
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "unauthorized.",
			})
			return
		}

		id, err := uuid.Parse(idString)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "unauthorized.",
			})
			return
		}

		user, err := userRepo.FindByID(id, true)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "internal server error.",
			})
			return
		}

		if user.UserGroup == nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "unauthorized.",
			})
			return
		}

		if !user.UserGroup.IsFullAccess {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "forbidden.",
			})
			return
		}

		ctx.Next()
	}
}
