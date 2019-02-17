package controller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/screwyprof/payment/internal/pkg/cqrs"
	"github.com/screwyprof/payment/internal/pkg/delivery/gin/responder"
	"github.com/screwyprof/payment/pkg/query"
	"github.com/screwyprof/payment/pkg/report"
)

// ShowAccount Retrieves account info.
type ShowAccount struct {
	queryBus cqrs.QueryHandler
}

// NewShowAccount Creates a new instance of ShowAccount.
func NewShowAccount(queryBus cqrs.QueryHandler) *ShowAccount {
	return &ShowAccount{queryBus: queryBus}
}

// Handle godoc
// @Summary Show an account
// @Description Show account info by number
// @Tags accounts
// @Accept  json
// @Produce  json
// @Param number path string true "account number"
// @Success 200 {object} response.AccountInfo
// @Failure 400 {object} responder.HTTPError
// @Failure 404 {object} responder.HTTPError
// @Failure 500 {object} responder.HTTPError
// @Router /accounts/{number} [get]
func (h *ShowAccount) Handle(ctx *gin.Context) {
	number := ctx.Param("number")

	accountReport := &report.Account{}
	err := h.queryBus.Handle(context.Background(), query.GetAccountShortInfo{Number: number}, accountReport)
	if err != nil {
		responder.NewError(ctx, http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusOK, AccountReportToAccountResponse(accountReport))
}
