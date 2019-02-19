package controller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/screwyprof/payment/internal/pkg/cqrs"
	"github.com/screwyprof/payment/internal/pkg/delivery/gin/response"

	"github.com/screwyprof/payment/pkg/query"
	"github.com/screwyprof/payment/pkg/report"
)

// ShowAvailableAccounts Retrieves available accounts.
type ShowAvailableAccounts struct {
	queryBus cqrs.QueryHandler
}

// NewShowAvailableAccounts Creates a new instance of ShowAvailableAccounts.
func NewShowAvailableAccounts(queryBus cqrs.QueryHandler) *ShowAvailableAccounts {
	return &ShowAvailableAccounts{queryBus: queryBus}
}

// Handle godoc
// @Summary Retrieves available accounts
// @Description Retrieves available accounts
// @Tags accounts
// @Accept  json
// @Produce  json
// @Success 200 {array} response.AvailableAccount
// @Failure 400 {object} response.HTTPError
// @Failure 404 {object} response.HTTPError
// @Failure 500 {object} response.HTTPError
// @Router /accounts [get]
func (h *ShowAvailableAccounts) Handle(ctx *gin.Context) {
	accs := report.Accounts{}
	err := h.queryBus.Handle(context.Background(), query.GetAllAccounts{}, &accs)
	if err != nil {
		response.NewError(ctx, http.StatusNotFound, err)
		return
	}

	var res []response.AvailableAccount
	for _, acc := range accs {
		res = append(res, response.AvailableAccount{
			Number: acc.Number,
		})
	}

	ctx.JSON(http.StatusOK, res)
}
