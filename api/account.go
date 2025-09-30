package api

import (
	"errors"
	"math/big"
	"net/http"

	db "example.com/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

func (server *Server) CreateAccount(c *gin.Context) {
	userId, ok := c.Get(authenticatedPayload)

	if !ok {
		c.JSON(http.StatusUnauthorized, errorResponse(errors.New("not authorized")))
		return

	}
	var payload db.CreateAccountParams

	if err := c.ShouldBindBodyWithJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if userId.(string) != payload.Owner {
		c.JSON(http.StatusBadRequest, errorResponse(errors.New("wrong request")))
		return
	}

	payload.Balance = pgtype.Numeric{
		Int:   big.NewInt(0),
		Exp:   0,
		Valid: true,
	}

	var account *db.Account

	err := server.Store.ExecTx(c, func(q *db.Queries) error {
		var err error


		exists, err := q.UserExists(c, payload.Owner)
		if err != nil{
			return err
		}
		if !exists {
			return errors.New("owner does not exist")
		}

		respAccount, err := q.CreateAccount(c, db.CreateAccountParams{
			Owner: payload.Owner,
			Currency: payload.Currency,
			Balance: payload.Balance,
		})

		if err != nil {
			return err
		}

		account = &respAccount
		return nil
	})

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.Code {
			case "23505", "23503":
				c.JSON(http.StatusForbidden, errorResponse(pgErr))
				return
			}
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, *account)

}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) GetAccount(c *gin.Context) {

	userId, ok := c.Get(authenticatedPayload)

	if !ok {
		c.JSON(http.StatusUnauthorized, errorResponse(errors.New("not authorized")))
		return

	}

	var req getAccountRequest

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.Store.GetAccountFromOwner(c, db.GetAccountFromOwnerParams{
		ID:    req.ID,
		Owner: userId.(string),
	})

	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusAccepted, account)
}

type listAccountsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) ListAccounts(c *gin.Context) {
	userId, ok := c.Get(authenticatedPayload)

	if !ok {
		c.JSON(http.StatusUnauthorized, errorResponse(errors.New("not authorized")))
		return

	}
	_, err := server.Store.GetUser(c, userId.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var req listAccountsRequest

	if err := c.BindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	accounts, err := server.Store.ListAccounts(c, db.ListAccountsParams{
		Owner:  userId.(string),
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusAccepted, accounts)
}
