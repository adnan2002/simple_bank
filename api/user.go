package api

import (
	"errors"
	"net/http"
	"time"

	db "example.com/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
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

	token, err := server.TokenMaker.CreateToken(user.Username, server.Config.AccessTokenDuration)

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := struct {
		Username      string        `json:"username"`
		FullName      string        `json:"full_name"`
		Email         string        `json:"email"`
		Token         string        `json:"token"`
		TokenDuration time.Duration `json:"token_duration"`
	}{
		Username:      user.Username,
		FullName:      user.FullName,
		Email:         user.Email,
		Token:         token,
		TokenDuration: server.Config.AccessTokenDuration,
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
	Username string `json:"username" binding:"required,alphanum"`
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
	token, err := server.TokenMaker.CreateToken(user.Username, server.Config.AccessTokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// 4. Return response
	c.JSON(http.StatusOK, gin.H{
		"username":       user.Username,
		"full_name":      user.FullName,
		"email":          user.Email,
		"token":          token,
		"token_duration": server.Config.AccessTokenDuration,
	})
}
