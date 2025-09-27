package api

import (
	"net/http"

	db "example.com/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)


type CreateUserRequest struct {
	Username     string `json:"username" binding:"required"`
    FullName     string `json:"full_name" binding:"required"`
    Email        string `json:"email" binding:"required"`
    Password 	 string `json:"password" binding:"required"`
}

func (server *Server) CreateUser(c *gin.Context) {
	var payload CreateUserRequest

	if err := c.ShouldBindBodyWithJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), 0)

	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.Store.CreateUser(c, db.CreateUserParams{
		Username:    payload.Username,
		FullName:  	payload.FullName,
		Email:   payload.Email,
		PasswordHash: string(hashedPassword),
	})

	if err != nil {
		if pgErr, ok  := err.(*pgconn.PgError); ok {
			switch pgErr.Code {
				case "23505":
				c.JSON(http.StatusForbidden, errorResponse(pgErr))
				return
			}
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, user)

}



type getUserRequest struct {
	Username string `uri:"username" binding:"required,min=1"`
}

func (server *Server) GetUser(c *gin.Context) {


	var req getUserRequest

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.Store.GetUser(c, req.Username)

	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusAccepted, user)
}
