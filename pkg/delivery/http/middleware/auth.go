package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)

		id := session.Get("id")
		if id == nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "unauthorized.",
			})
		}
		ctx.Next()
	}
}

func AuthIsFullAccess() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)

		id := session.Get("id")
		isFullAccess := session.Get("isFullAccess")
		if id == nil || isFullAccess == nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "unauthorized.",
			})
		}

		if is, ok := isFullAccess.(bool); !ok || !is {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "forbidden.",
			})
		}

		ctx.Next()
	}
}
