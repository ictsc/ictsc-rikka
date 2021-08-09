package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ictsc/ictsc-rikka/pkg/error"
)

type ErrorMiddleware struct {
}

func NewErrorMiddleware() *ErrorMiddleware {
	return &ErrorMiddleware{}
}

func (m *ErrorMiddleware) HandleError(ctx *gin.Context) {
	ctx.Next()

	if len(ctx.Errors) != 0 {
		var statusCode int
		err := ctx.Errors.Last().Err
		switch err.(type) {
		case *error.BadRequestError:
			statusCode = http.StatusBadRequest
		case *error.UnauthorizedError:
			statusCode = http.StatusUnauthorized
		case *error.ForbiddenError:
			statusCode = http.StatusForbidden
		case *error.NotFoundError:
			statusCode = http.StatusNotFound
		case *error.InternalServerError:
			statusCode = http.StatusInternalServerError
		default:
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "unknown error",
			})
		}
		ctx.AbortWithStatusJSON(statusCode, gin.H{"error": err.Error()})
	}

}
