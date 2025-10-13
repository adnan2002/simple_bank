package api

import (
	"net/http"
	"time"

	db "example.com/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}


func (server *Server) Refresh(c *gin.Context) {
	var refresh_token RefreshRequest
	if err := c.ShouldBindJSON(&refresh_token); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, err := server.TokenMaker.VerifyToken(refresh_token.RefreshToken)

	if err != nil {
		c.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}


	accessToken, _, err := server.TokenMaker.CreateToken(payload.Username, server.Config.AccessTokenDuration)

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, refreshPayload, err := server.TokenMaker.CreateToken(payload.Username, server.Config.RefreshTokenduration)

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	_, err = server.Store.CreateSession(c, db.CreateSessionParams{
		ID: pgtype.UUID{
			Bytes: [16]byte(refreshPayload.ID.Bytes()),
			Valid: true,
		},
		Username:     payload.Username,
		RefreshToken: refreshToken,
		UserAgent:    c.Request.UserAgent(),
		ClientIp:     c.ClientIP(),
		ExpiresAt: pgtype.Timestamptz{
			Time:             refreshPayload.ExpiredAt,
			InfinityModifier: pgtype.Finite,
			Valid:            true,
		},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

		c.SetCookie(
		"refresh_token",                                    // name
		refreshToken,                                       // value
		int(server.Config.RefreshTokenduration.Seconds()), // maxAge in seconds
		"/",                                                // path
		"",                                                 // domain (empty means current domain)
		true,                                               // secure (true for HTTPS only)
		true,                                               // httpOnly
	)

	response := struct {
		AccessToken         string        `json:"access_token"`
		AccessTokenDuration time.Duration `json:"access_token_duration"`
	}{
		AccessToken:         accessToken,
		AccessTokenDuration: server.Config.AccessTokenDuration,
	}


	c.JSON(http.StatusOK, response)






}