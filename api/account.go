package api

import (
	"math/big"
	"net/http"

	db "example.com/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (server *Server) CreateAccount(c *gin.Context) {
	var payload db.CreateAccountParams

	if err := c.ShouldBindBodyWithJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload.Balance = pgtype.Numeric{
		Int:   big.NewInt(0),
		Exp:   0,
		Valid: true,
	}

	account, err := server.Store.CreateAccount(c, db.CreateAccountParams{
		Owner:    payload.Owner,
		Currency: payload.Currency,
		Balance:  payload.Balance,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, account)

}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) GetAccount(c *gin.Context) {
	// id, _ := c.Params.Get("id")

	// idInt, _ := strconv.ParseInt(id, 10, 64)

	// account, err := server.Store.GetAccount(c, int64(idInt))

	// if err != nil {
	// 	c.JSON(http.StatusNotFound, errorResponse(err))
	// 	return
	// }

	// c.JSON(http.StatusAccepted, account)

	var req getAccountRequest

	if err := c.ShouldBindUri(&req); err!= nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.Store.GetAccount(c, req.ID)

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
	var req listAccountsRequest

	if err := c.BindQuery(&req); err!= nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	accounts, err := server.Store.ListAccounts(c, db.ListAccountsParams{
		Limit: req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	})

	if err != nil {

		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusAccepted, accounts)
}


