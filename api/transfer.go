package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/hamdysherif/simplebank/db/sqlc"
)

type transferRequest struct {
	FromAccount int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccount   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount      int64  `json:"amount" binding:"required,gt=0"`
	Currency    string `json:"currency" binding:"required,currency"`
}

func (server *Server) transferAmount(ctx *gin.Context) {
	var req transferRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, responseError(err))
		return
	}

	if !isValidAccount(server, ctx, req.FromAccount, req.Currency) {
		return
	}
	if !isValidAccount(server, ctx, req.ToAccount, req.Currency) {
		return
	}

	result, err := server.db.TransferTx(ctx, db.TransferParams{FromAccountID: req.FromAccount, ToAccountID: req.ToAccount, Amount: req.Amount})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, responseError(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func isValidAccount(server *Server, ctx *gin.Context, accountID int64, currency string) bool {
	account, err := server.db.GetAccount(ctx, accountID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, responseError(err))
		return false
	}
	if account.Currency != currency {
		ctx.JSON(http.StatusBadRequest, responseError(fmt.Errorf("account [%v], currency not match %v vs %v", account.ID, account.Currency, currency)))
		return false
	}
	return true
}
