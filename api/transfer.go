package api

import (
	"context"
	"net/http"

	db "example.com/db/sqlc"
	"github.com/gin-gonic/gin"
)



type RequestParams struct {
	FromAccountId int64  `json:"from_account_id"`
	ToAccountId   int64  `json:"to_account_id"`
	Amount        int64  `json:"amount"`
	Currency      string `json:"currency" binding:"currency,required"`
}

func (server *Server) CreateTransfer(c *gin.Context) {
	var payload RequestParams
	if err := c.ShouldBindBodyWithJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ok := server.IsSameCurrency(payload.ToAccountId, payload.Currency)

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Currencies don't match yo"})
		return
	}

	transfer, err := server.Store.TransferTx(c, db.TransferTxParams{
		FromAccountId: payload.FromAccountId,
		ToAccountId: payload.ToAccountId,
		Amount: payload.Amount,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, transfer)

}



func (server *Server) IsSameCurrency(toAccountId int64, currency string) bool {
	account, err := server.Store.GetAccount(context.Background(), toAccountId)

	if err != nil {
		return false
	}

	if account.Currency != currency {
		return false
	}

	return true


}	




