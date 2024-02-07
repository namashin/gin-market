package middleware

import (
	"gin-market/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AuthMiddleware(authService services.IAuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.GetHeader("Authorization")

		if len(header) == 0 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		} else if !strings.HasPrefix(header, "Bearer ") {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(header, "Bearer ")
		user, err := authService.GetUserFromToken(tokenString)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.Set("user", user)

		// next flow
		ctx.Next()
	}
}
