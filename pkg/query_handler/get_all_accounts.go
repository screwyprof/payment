package query_handler

import (
	"context"
	"fmt"
	"github.com/screwyprof/payment/internal/pkg/cqrs"

	"github.com/screwyprof/payment/pkg/query"
	"github.com/screwyprof/payment/pkg/report"
)

type getAllAccountsInfo struct {
	accountProvider report.GetAllAccounts
}

func NewGetAllAccounts(accountProvider report.GetAllAccounts) cqrs.QueryHandler {
	return &getAllAccountsInfo{accountProvider: accountProvider}
}

func (intr *getAllAccountsInfo) Handle(ctx context.Context, q cqrs.Query, result interface{}) error {
	_, ok := q.(query.GetAllAccounts)
	if !ok {
		return fmt.Errorf("invalid query %+#v given", q)
	}

	res, ok := result.(*report.Accounts)
	if !ok {
		return fmt.Errorf("invalid report %+#v given", result)
	}

	accs, err := intr.accountProvider.All()
	if err != nil {
		return fmt.Errorf("cannot retrieve accounts: %v", err)
	}

	for _, acc := range accs {
		*res = append(*res, acc)
	}

	return nil
}
