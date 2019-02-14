package query_handler

import (
	"context"
	"fmt"

	"github.com/screwyprof/payment/internal/pkg/cqrs"

	"github.com/screwyprof/payment/pkg/query"
	"github.com/screwyprof/payment/pkg/report"
)

type getAccountShortInfo struct {
	accountProvider report.GetAccountByNumber
}

func NewGetAccountShortInfo(accountProvider report.GetAccountByNumber) cqrs.QueryHandler {
	return &getAccountShortInfo{accountProvider: accountProvider}
}

func (intr *getAccountShortInfo) Handle(ctx context.Context, q cqrs.Query, result interface{}) error {
	req, ok := q.(query.GetAccountShortInfo)
	if !ok {
		return fmt.Errorf("invalid query %+#v given", q)
	}

	res, ok := result.(*report.Account)
	if !ok {
		return fmt.Errorf("invalid report %+#v given", result)
	}

	acc, err := intr.accountProvider.ByNumber(req.Number)
	if err != nil {
		return fmt.Errorf("cannot retrieve account: %v", err)
	}

	*(res) = *acc

	return nil
}
