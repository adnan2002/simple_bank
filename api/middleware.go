package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"example.com/token"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderString = "Authorization"
	authorizationPrefixString = "Bearer"
	authenticatedPayload      = "auth_username"
)

func AuthMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(authorizationHeaderString)
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(errors.New("authorization header is missing")))
			return
		}

		bearerPrefix := fmt.Sprintf("%s ", authorizationPrefixString)
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(errors.New("invalid authorization header format")))
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, bearerPrefix)

		payload, err := tokenMaker.VerifyToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		fmt.Println("payload username is :",payload.Username)

		c.Set(authenticatedPayload, payload.Username)

		c.Next()
	}
}
