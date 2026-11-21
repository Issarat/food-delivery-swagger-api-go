package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/Issarat/food-delivery-swagger-api-go/api-go/internal/auth"
)

func RequireJWT(tokens *auth.TokenManager, service *auth.Service) gin.HandlerFunc {	return func(c *gin.Context) {
		parts := strings.Fields(c.GetHeader("Authorization"))
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.Header("WWW-Authenticate", "Bearer")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "bearer token is required",
			})
			return
		}

		claims, err := tokens.Verify(parts[1])
		if err != nil {
			c.Header("WWW-Authenticate", `Bearer error="invalid_token"`)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			return
		}
		if claims.ID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "token ID is required",
			})
			return
		}

		revoked, err := service.IsTokenRevoked(c.Request.Context(), claims.ID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
				"error": "could not check token status",
			})
			return
		}

		if revoked {
			c.Header("WWW-Authenticate", `Bearer error="invalid_token"`)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "token has been revoked",
			})
			return
		}

		c.Set("user_id", claims.Subject)
		c.Set("token_jti", claims.ID)
		c.Set("token_expires_at", claims.ExpiresAt.Time)
		c.Next()
	}
}
