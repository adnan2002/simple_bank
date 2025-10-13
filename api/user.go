package api

import (
	"errors"
	"net/http"
	"time"

	db "example.com/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (server *Server) CreateUser(c *gin.Context) {
	var payload CreateUserRequest

	if err := c.ShouldBindBodyWithJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)

	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.Store.CreateUser(c, db.CreateUserParams{
		Username:     payload.Username,
		FullName:     payload.FullName,
		Email:        payload.Email,
		PasswordHash: string(hashedPassword),
	})

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.Code {
			case "23505":
				c.JSON(http.StatusForbidden, errorResponse(pgErr))
				return
			}
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	accessToken, _, err := server.TokenMaker.CreateToken(user.Username, server.Config.AccessTokenDuration)

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, refreshAccessPayload, err := server.TokenMaker.CreateToken(user.Username, server.Config.RefreshTokenduration)

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	_, err = server.Store.CreateSession(c, db.CreateSessionParams{
		ID: pgtype.UUID{
			Bytes: [16]byte(refreshAccessPayload.ID.Bytes()),
			Valid: true,
		},
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    c.Request.UserAgent(),
		ClientIp:     c.ClientIP(),
		ExpiresAt: pgtype.Timestamptz{
			Time:             refreshAccessPayload.ExpiredAt,
			InfinityModifier: pgtype.Finite,
			Valid:            true,
		},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Set refresh token as HTTP-only cookie
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
		SessionID           uuid.UUID     `json:"session_id"`
		AccessToken         string        `json:"access_token"`
		AccessTokenDuration time.Duration `json:"access_token_duration"`
		Username            string        `json:"username"`
		FullName            string        `json:"full_name"`
		Email               string        `json:"email"`
	}{
		SessionID:           refreshAccessPayload.ID,
		AccessToken:         accessToken,
		AccessTokenDuration: server.Config.AccessTokenDuration,
		Username:            user.Username,
		FullName:            user.FullName,
		Email:               user.Email,
	}

	c.JSON(http.StatusCreated, response)
}

func (server *Server) GetUser(c *gin.Context) {

	userId, ok := c.Get(authenticatedPayload)

	if !ok {
		c.JSON(http.StatusUnauthorized, errorResponse(errors.New("not authorized")))
		return
	}
	
	// 4. Use payload.Username to fetch the user
	user, err := server.Store.GetUser(c, userId.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"username":  user.Username,
		"full_name": user.FullName,
		"email":     user.Email,
	})
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

func (server *Server) LoginUser(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// 1. Fetch user from database
	user, err := server.Store.GetUser(c, req.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errorResponse(errors.New("invalid username or password")))
		return
	}

	// 2. Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, errorResponse(errors.New("invalid username or password")))
		return
	}

	// 3. Create access token
	accessToken, _, err := server.TokenMaker.CreateToken(user.Username, server.Config.AccessTokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// 4. Create refresh token
	refreshToken, refreshAccessPayload, err := server.TokenMaker.CreateToken(user.Username, server.Config.RefreshTokenduration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// 5. Create session
	_, err = server.Store.CreateSession(c, db.CreateSessionParams{
		ID: pgtype.UUID{
			Bytes: [16]byte(refreshAccessPayload.ID.Bytes()),
			Valid: true,
		},
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    c.Request.UserAgent(),
		ClientIp:     c.ClientIP(),
		ExpiresAt: pgtype.Timestamptz{
			Time:             refreshAccessPayload.ExpiredAt,
			InfinityModifier: pgtype.Finite,
			Valid:            true,
		},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Set refresh token as HTTP-only cookie
	c.SetCookie(
		"refresh_token",                                    // name
		refreshToken,                                       // value
		int(server.Config.RefreshTokenduration.Seconds()), // maxAge in seconds
		"/",                                                // path
		"",                                                 // domain (empty means current domain)
		true,                                               // secure (true for HTTPS only)
		true,                                               // httpOnly
	)

	// 6. Return response (without refresh token)
	response := struct {
		SessionID           uuid.UUID     `json:"session_id"`
		AccessToken         string        `json:"access_token"`
		AccessTokenDuration time.Duration `json:"access_token_duration"`
		Username            string        `json:"username"`
		FullName            string        `json:"full_name"`
		Email               string        `json:"email"`
	}{
		SessionID:           refreshAccessPayload.ID,
		AccessToken:         accessToken,
		AccessTokenDuration: server.Config.AccessTokenDuration,
		Username:            user.Username,
		FullName:            user.FullName,
		Email:               user.Email,
	}

	c.JSON(http.StatusOK, response)
}