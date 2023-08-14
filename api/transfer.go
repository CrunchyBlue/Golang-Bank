package api

import (
	"database/sql"
	"errors"
	"fmt"
	db "github.com/CrunchyBlue/Golang-Bank/sqlc"
	"github.com/CrunchyBlue/Golang-Bank/token"
	"github.com/gin-gonic/gin"
	"net/http"
)

type getTransfersRequest struct {
	PageNumber int32 `form:"page_number" binding:"required,min=1"`
	PageSize   int32 `form:"page_size" binding:"required,min=10,max=50"`
}

func (server *Server) getTransfers(ctx *gin.Context) {
	var req getTransfersRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.GetTransfersParams{
		Limit:  req.PageSize,
		Offset: (req.PageNumber - 1) * req.PageSize,
	}

	transfers, err := server.store.GetTransfers(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, transfers)
}

type getOutboundTransfersForAccountUriParams struct {
	AccountID int64 `uri:"account_id" binding:"required,min=1"`
}

type getOutboundTransfersForAccountQueryParams struct {
	PageNumber int32 `form:"page_number" binding:"required,min=1"`
	PageSize   int32 `form:"page_size" binding:"required,min=10,max=50"`
}

type getOutboundTransfersForAccountRequest struct {
	UriParams   getOutboundTransfersForAccountUriParams
	QueryParams getOutboundTransfersForAccountQueryParams
}

func (server *Server) getOutboundTransfersForAccount(ctx *gin.Context) {
	var req getOutboundTransfersForAccountRequest

	if err := ctx.ShouldBindUri(&req.UriParams); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindQuery(&req.QueryParams); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.GetOutboundTransfersForAccountParams{
		SourceAccountID: req.UriParams.AccountID,
		Limit:           req.QueryParams.PageSize,
		Offset:          (req.QueryParams.PageNumber - 1) * req.QueryParams.PageSize,
	}

	entries, err := server.store.GetOutboundTransfersForAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, entries)
}

type getInboundTransfersForAccountUriParams struct {
	AccountID int64 `uri:"account_id" binding:"required,min=1"`
}

type getInboundTransfersForAccountQueryParams struct {
	PageNumber int32 `form:"page_number" binding:"required,min=1"`
	PageSize   int32 `form:"page_size" binding:"required,min=10,max=50"`
}

type getInboundTransfersForAccountRequest struct {
	UriParams   getInboundTransfersForAccountUriParams
	QueryParams getInboundTransfersForAccountQueryParams
}

func (server *Server) getInboundTransfersForAccount(ctx *gin.Context) {
	var req getInboundTransfersForAccountRequest

	if err := ctx.ShouldBindUri(&req.UriParams); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindQuery(&req.QueryParams); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.GetInboundTransfersForAccountParams{
		DestinationAccountID: req.UriParams.AccountID,
		Limit:                req.QueryParams.PageSize,
		Offset:               (req.QueryParams.PageNumber - 1) * req.QueryParams.PageSize,
	}

	entries, err := server.store.GetInboundTransfersForAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, entries)
}

type getTransferRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getTransfer(ctx *gin.Context) {
	var req getTransferRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	transfer, err := server.store.GetTransfer(ctx, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, transfer)
}

type createTransferRequest struct {
	SourceAccountID      int64  `json:"source_account_id" binding:"required,min=1"`
	DestinationAccountID int64  `json:"destination_account_id" binding:"required,min=1"`
	Amount               int64  `json:"amount" binding:"required,min=1"`
	Currency             string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req createTransferRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	sourceAccount, valid := server.validateAccount(ctx, req.SourceAccountID, req.Currency)
	if !valid {
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if sourceAccount.Owner != authPayload.Username {
		err := errors.New("source account does not belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	_, valid = server.validateAccount(ctx, req.DestinationAccountID, req.Currency)
	if !valid {
		return
	}

	arg := db.TransferTxParams{
		SourceAccountID:      req.SourceAccountID,
		DestinationAccountID: req.DestinationAccountID,
		Amount:               req.Amount,
	}

	transfer, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, transfer)
}

type updateTransferUriParams struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type updateTransferBody struct {
	Amount int64 `json:"amount" binding:"required,min=1"`
}

type updateTransferRequest struct {
	UriParams updateTransferUriParams
	Body      updateTransferBody
}

func (server *Server) updateTransfer(ctx *gin.Context) {
	var req updateTransferRequest

	if err := ctx.ShouldBindUri(&req.UriParams); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindJSON(&req.Body); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateTransferParams{
		ID:     req.UriParams.ID,
		Amount: req.Body.Amount,
	}

	transfer, err := server.store.UpdateTransfer(ctx, arg)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, transfer)
}

type deleteTransferRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteTransfer(ctx *gin.Context) {
	var req deleteTransferRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteTransfer(ctx, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.Status(http.StatusOK)
}

func (server *Server) validateAccount(ctx *gin.Context, accountID int64, currency string) (db.Account, bool) {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return account, false
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", account.ID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false
	}

	return account, true
}
