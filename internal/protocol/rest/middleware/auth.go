package middleware

import (
	"diary-api/internal/auth"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func JwtMiddleware(tokenService auth.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtToken := extractJwt(c)
		if jwtToken == "" {
			returnUnauthorized(c)
			return
		}
		claims, err := tokenService.ValidateToken(jwtToken)
		if err != nil {
			returnUnauthorized(c)
			return
		}

		c.Set(auth.AuthToken, jwtToken)
		c.Set(auth.UserIDKey, claims.UserID)
		c.Next()
	}
}

func extractJwt(ctx *gin.Context) string {
	jwtToken := ctx.GetHeader("Authorization")
	if jwtToken == "" {
		return ""
	}
	splitToken := strings.Split(jwtToken, "Bearer ")
	if len(splitToken) != 2 {
		return ""
	}
	jwtToken = splitToken[1]
	return jwtToken
}

func returnUnauthorized(ctx *gin.Context) {
	ctx.JSON(http.StatusUnauthorized, gin.H{"message": "jwt token is missing or invalid"})
	ctx.Abort()
}
