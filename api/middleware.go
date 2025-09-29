package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"example.com/token"
)

// AuthMiddleware creates a gin.HandlerFunc that checks JWT/Paseto token validity
func AuthMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Extract the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(errors.New("authorization header is missing")))
			return
		}

		// 2. Check if it is a Bearer token
		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(errors.New("invalid authorization header format")))
			return
		}

		// 3. Get the token string
		tokenStr := strings.TrimPrefix(authHeader, bearerPrefix)

		// 4. Verify token
		payload, err := tokenMaker.VerifyToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		// 5. Save the payload in Gin context (so handlers can use it)
		c.Set("auth_username", payload.Username)

		// 6. Continue to the next handler
		c.Next()
	}
}






