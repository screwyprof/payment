package controller

import (
	"context"
	"github.com/google/uuid"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rhymond/go-money"

	"github.com/screwyprof/payment/internal/pkg/delivery/gin/request"
	"github.com/screwyprof/payment/internal/pkg/delivery/gin/response"

	"github.com/screwyprof/payment/pkg/command"
	"github.com/screwyprof/payment/pkg/domain"
	"github.com/screwyprof/payment/pkg/query"
	"github.com/screwyprof/payment/pkg/report"
)

type OpenAccount struct {
	commandBus domain.CommandHandler
	queryBus   domain.QueryHandler
}

func NewOpenAccount(commandBus domain.CommandHandler, queryBus domain.QueryHandler) *OpenAccount {
	return &OpenAccount{
		commandBus: commandBus,
		queryBus:   queryBus,
	}
}

// Handle godoc
// @Summary Open a new account
// @Description Open a new account with optional balance
// @Tags accounts
// @Accept  json
// @Produce  json
// @Param account body request.OpenAccount true "Open account"
// @Success 200 {object} response.ShortAccountInfo
// @Failure 400 {object} response.HTTPError
// @Failure 404 {object} response.HTTPError
// @Failure 500 {object} response.HTTPError
// @Router /accounts [post]
func (h *OpenAccount) Handle(ctx *gin.Context) {
	var req request.OpenAccount
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	//if err := req.Validate(); err != nil {
	//	response.NewError(ctx, http.StatusBadRequest, err)
	//	return
	//}

	err := h.commandBus.Handle(context.Background(), command.OpenAccount{
		AggID:   uuid.New(),
		Number:  req.Number,
		Balance: *money.New(req.Amount, req.Currency),
	})
	if err != nil {
		response.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	accountReport := &report.Account{}
	err = h.queryBus.Handle(context.Background(), query.GetAccountShortInfo{Number: req.Number}, accountReport)
	if err != nil {
		response.NewError(ctx, http.StatusNotFound, err)
		return
	}

	resp := response.ShortAccountInfo{Number: accountReport.Number, Balance: accountReport.Balance.Display()}
	ctx.JSON(http.StatusOK, resp)
}
