package controller

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rhymond/go-money"

	"github.com/screwyprof/payment/internal/pkg/cqrs"

	"github.com/screwyprof/payment/internal/pkg/delivery/gin/request"
	"github.com/screwyprof/payment/internal/pkg/delivery/gin/responder"
	"github.com/screwyprof/payment/pkg/command"
	"github.com/screwyprof/payment/pkg/domain/account"
)

// TransferMoney Tansfers money from one account to another.
type TransferMoney struct {
	commandBus cqrs.CommandHandler
	queryBus   cqrs.QueryHandler
}

// NewTransferMoney Creates a new instance of TransferMoney.
func NewTransferMoney(commandBus cqrs.CommandHandler, queryBus cqrs.QueryHandler) *TransferMoney {
	return &TransferMoney{
		commandBus: commandBus,
		queryBus:   queryBus,
	}
}

// Handle godoc
// @Summary Transfer money from an account to another account
// @Description Transfer money from an account to another account
// @Tags accounts
// @Accept  json
// @Produce  json
// @Param number path string true "account number"
// @Param transfer body request.TransferMoney true "Transfer money"
// @Success 200 {object} controller.Message
// @Failure 400 {object} responder.HTTPError
// @Failure 500 {object} responder.HTTPError
// @Router /accounts/{number}/transfer [post]
func (h *TransferMoney) Handle(ctx *gin.Context) {
	from := ctx.Param("number")

	var req request.TransferMoney
	if err := ctx.ShouldBindJSON(&req); err != nil {
		responder.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	//if err := req.Validate(); err != nil {
	//	responder.NewError(ctx, http.StatusBadRequest, err)
	//	return
	//}

	amount := *money.New(req.Amount, req.Currency)

	err := h.commandBus.Handle(context.Background(), command.TransferMoney{
		From:   account.Number(from),
		To:     account.Number(req.To),
		Amount: amount,
	})
	if err != nil {
		responder.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	err = h.commandBus.Handle(context.Background(), command.ReceiveMoney{
		From:   account.Number(from),
		To:     account.Number(req.To),
		Amount: amount,
	})
	if err != nil {
		responder.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, Message{
		Message: fmt.Sprintf("%s transfered from %s to %s", amount.Display(), from, req.To)})
}
