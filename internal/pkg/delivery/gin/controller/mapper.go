package controller

import (
	"github.com/screwyprof/payment/internal/pkg/delivery/gin/response"

	"github.com/screwyprof/payment/pkg/report"
)

func AccountReportToAccountResponse(accountReport *report.Account) response.AccountInfo {
	var ledgers []response.Ledger
	for _, ledger := range accountReport.Ledgers {
		ledgers = append(ledgers, response.Ledger{Action: ledger.ToString()})
	}

	return response.AccountInfo{
		Number:  accountReport.Number,
		Balance: accountReport.Balance.Display(),
		Ledgers: ledgers,
	}
}
