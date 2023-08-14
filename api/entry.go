package api

import (
	"database/sql"
	"errors"
	db "github.com/CrunchyBlue/Golang-Bank/sqlc"
	"github.com/gin-gonic/gin"
	"net/http"
)

type getEntriesRequest struct {
	PageNumber int32 `form:"page_number" binding:"required,min=1"`
	PageSize   int32 `form:"page_size" binding:"required,min=10,max=50"`
}

func (server *Server) getEntries(ctx *gin.Context) {
	var req getEntriesRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.GetEntriesParams{
		Limit:  req.PageSize,
		Offset: (req.PageNumber - 1) * req.PageSize,
	}

	entries, err := server.store.GetEntries(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, entries)
}

type getEntriesForAccountUriParams struct {
	AccountID int64 `uri:"account_id" binding:"required,min=1"`
}

type getEntriesForAccountQueryParams struct {
	PageNumber int32 `form:"page_number" binding:"required,min=1"`
	PageSize   int32 `form:"page_size" binding:"required,min=10,max=50"`
}

type getEntriesForAccountRequest struct {
	UriParams   getEntriesForAccountUriParams
	QueryParams getEntriesForAccountQueryParams
}

func (server *Server) getEntriesForAccount(ctx *gin.Context) {
	var req getEntriesForAccountRequest

	if err := ctx.ShouldBindUri(&req.UriParams); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindQuery(&req.QueryParams); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.GetEntriesForAccountParams{
		AccountID: req.UriParams.AccountID,
		Limit:     req.QueryParams.PageSize,
		Offset:    (req.QueryParams.PageNumber - 1) * req.QueryParams.PageSize,
	}

	entries, err := server.store.GetEntriesForAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, entries)
}

type getEntryRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getEntry(ctx *gin.Context) {
	var req getEntryRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	entry, err := server.store.GetEntry(ctx, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, entry)
}

type createEntryRequest struct {
	AccountID int64 `json:"account_id" binding:"required,min=1"`
	Amount    int64 `json:"amount" binding:"required"`
}

func (server *Server) createEntry(ctx *gin.Context) {
	var req createEntryRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateEntryParams{
		AccountID: req.AccountID,
		Amount:    req.Amount,
	}

	entry, err := server.store.CreateEntry(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, entry)
}

type updateEntryUriParams struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type updateEntryBody struct {
	Amount int64 `json:"amount" binding:"required"`
}

type updateEntryRequest struct {
	UriParams updateEntryUriParams
	Body      updateEntryBody
}

func (server *Server) updateEntry(ctx *gin.Context) {
	var req updateEntryRequest

	if err := ctx.ShouldBindUri(&req.UriParams); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindJSON(&req.Body); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateEntryParams{
		ID:     req.UriParams.ID,
		Amount: req.Body.Amount,
	}

	entry, err := server.store.UpdateEntry(ctx, arg)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, entry)
}

type deleteEntryRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteEntry(ctx *gin.Context) {
	var req deleteEntryRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteEntry(ctx, req.ID)
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
