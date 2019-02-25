package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/rhymond/go-money"

	qdaptor "github.com/screwyprof/payment/internal/pkg/adaptor/query_handler"
	"github.com/screwyprof/payment/internal/pkg/bus"
	"github.com/screwyprof/payment/internal/pkg/cqrs"
	"github.com/screwyprof/payment/internal/pkg/observer"
	"github.com/screwyprof/payment/internal/pkg/reporting"

	"github.com/screwyprof/payment/pkg/command"
	"github.com/screwyprof/payment/pkg/domain"
	"github.com/screwyprof/payment/pkg/domain/account"
	"github.com/screwyprof/payment/pkg/event_handler"
	"github.com/screwyprof/payment/pkg/query"
	"github.com/screwyprof/payment/pkg/query_handler"
	"github.com/screwyprof/payment/pkg/report"
)

func main() {
	// init deps
	cqrs.RegisterAggregate(func(id uuid.UUID) domain.Aggregate {
		return account.Create(id)
	})

	accountReporter := reporting.NewInMemoryAccountReporter()

	notifier := observer.NewNotifier()
	notifier.Register(event_handler.NewAccountOpened(accountReporter))
	notifier.Register(event_handler.NewMoneyTransfered(accountReporter))
	notifier.Register(event_handler.NewMoneyReceived(accountReporter))

	eventStore := cqrs.NewInMemoryEventStore()
	commandHandler := cqrs.NewEventSourceHandler(eventStore, notifier)

	accountQueryer := query_handler.NewGetAccountShortInfo(accountReporter)
	queryBus := newQueryBus(accountQueryer)

	// Run test script

	// create a couple of accounts
	acc1ID := uuid.New()
	acc1Number := "ACC500" //account.GenerateAccNumber()
	openAccount(commandHandler, acc1ID, acc1Number, *money.New(50000, "USD"))

	acc2ID := uuid.New()
	acc2Number := "ACC300" //account.GenerateAccNumber()
	openAccount(commandHandler, acc2ID, acc2Number, *money.New(30000, "USD"))

	// get accounts info
	queryAccount(queryBus, string(acc1Number))
	queryAccount(queryBus, string(acc2Number))

	// transfer money and show results
	transferMoney(commandHandler, acc1ID, acc2ID, acc1Number, acc2Number, *money.New(30000, "USD"))
	queryAccount(queryBus, string(acc1Number))
	queryAccount(queryBus, string(acc2Number))
}

func newQueryBus(accountQueryer domain.QueryHandler) domain.QueryHandler {
	queryBus := bus.NewQueryHandlerBus()
	queryBus.Register("GetAccountShortInfo", qdaptor.FromDomain(accountQueryer))

	return qdaptor.ToDomain(queryBus)
}

func openAccount(commandBus domain.CommandHandler, ID uuid.UUID, number string, balance money.Money) {
	openAccountCmd := command.OpenAccount{
		AggID:   ID,
		Number:  number,
		Balance: balance,
	}
	failOnError(commandBus.Handle(context.Background(), openAccountCmd))
}

func transferMoney(commandBus domain.CommandHandler, fromID, toID uuid.UUID, from, to string, amount money.Money) {
	transferCmd := command.TransferMoney{
		AggID:  fromID,
		From:   from,
		To:     to,
		Amount: amount,
	}
	failOnError(commandBus.Handle(context.Background(), transferCmd))

	receiveCmd := command.ReceiveMoney{
		AggID:  toID,
		From:   from,
		To:     to,
		Amount: amount,
	}
	failOnError(commandBus.Handle(context.Background(), receiveCmd))
}

func queryAccount(queryBus domain.QueryHandler, number string) {
	acc := &report.Account{}
	err := queryBus.Handle(context.Background(), query.GetAccountShortInfo{Number: number}, acc)
	failOnError(err)
	fmt.Println(acc.ToString())

	if len(acc.Ledgers) < 1 {
		return
	}
	fmt.Println("Ledgers:")
	for idx, ledger := range acc.Ledgers {
		fmt.Printf("#%d. %s\n", idx+1, ledger.ToString())
	}
}

func failOnError(err error) {
	if err != nil {
		fmt.Printf("an error occured: %v", err)
		os.Exit(1)
	}
}
