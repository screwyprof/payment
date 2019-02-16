package main

import (
	"context"
	"fmt"
	"os"

	"github.com/rhymond/go-money"

	cmdadaptor "github.com/screwyprof/payment/internal/pkg/adaptor/command_handler"
	qdaptor "github.com/screwyprof/payment/internal/pkg/adaptor/query_handler"
	"github.com/screwyprof/payment/internal/pkg/bus"
	"github.com/screwyprof/payment/internal/pkg/cqrs"
	"github.com/screwyprof/payment/internal/pkg/observer"
	"github.com/screwyprof/payment/internal/pkg/reporting"
	"github.com/screwyprof/payment/internal/pkg/repository"

	"github.com/screwyprof/payment/pkg/command"
	"github.com/screwyprof/payment/pkg/command_handler"
	"github.com/screwyprof/payment/pkg/domain/account"
	"github.com/screwyprof/payment/pkg/event_handler"
	"github.com/screwyprof/payment/pkg/query"
	"github.com/screwyprof/payment/pkg/query_handler"
	"github.com/screwyprof/payment/pkg/report"
)

func main() {
	// init deps
	accountReporter := reporting.NewInMemoryAccountReporter()
	accountRepo := repository.NewInMemoryAccountReporter()

	notifier := observer.NewNotifier()
	notifier.Register(event_handler.NewAccountOpened(accountReporter))
	notifier.Register(event_handler.NewMoneyTransfered(accountReporter))
	//notifier.Register(event_handler.NewMoneySentFurther(receiveMoney))
	notifier.Register(event_handler.NewMoneyReceived(accountReporter))

	accountOpenner := command_handler.NewOpenAccount(accountRepo, notifier)
	moneyTransfer := command_handler.NewTransferMoney(accountRepo, accountRepo, notifier)
	moneyReceiver := command_handler.NewReceiveMoney(accountRepo, accountRepo, notifier)
	commandBus := newCommandBus(moneyTransfer, moneyReceiver, accountOpenner)

	accountQueryer := query_handler.NewGetAccountShortInfo(accountReporter)
	queryBus := newQueryBus(accountQueryer)

	// Run test script

	// create a couple of accounts
	acc1Number := account.Number("ACC500") //account.GenerateAccNumber()
	openAccount(commandBus, acc1Number, *money.New(50000, "USD"))

	acc2Number := account.Number("ACC300") //account.GenerateAccNumber()
	openAccount(commandBus, acc2Number, *money.New(30000, "USD"))

	// get accounts info
	queryAccount(queryBus, string(acc1Number))
	queryAccount(queryBus, string(acc2Number))

	// transfer money and show results
	transferMoney(commandBus, acc1Number, acc2Number, *money.New(30000, "USD"))
	queryAccount(queryBus, string(acc1Number))
	queryAccount(queryBus, string(acc2Number))
}

func newCommandBus(moneyTransfer, moneyReceiver, accountOpenner cqrs.CommandHandler) cqrs.CommandHandler {
	commandBus := bus.NewCommandBus()
	commandBus.Register("OpenAccount", cmdadaptor.FromDomain(accountOpenner))
	commandBus.Register("TransferMoney", cmdadaptor.FromDomain(moneyTransfer))
	commandBus.Register("ReceiveMoney", cmdadaptor.FromDomain(moneyReceiver))

	return cmdadaptor.ToDomain(commandBus)
}

func newQueryBus(accountQueryer cqrs.QueryHandler) cqrs.QueryHandler {
	queryBus := bus.NewQueryHandlerBus()
	queryBus.Register("GetAccountShortInfo", qdaptor.FromDomain(accountQueryer))

	return qdaptor.ToDomain(queryBus)
}

func openAccount(commandBus cqrs.CommandHandler, number account.Number, balance money.Money) {
	openAccountCmd := command.OpenAccount{
		Number:  number,
		Balance: balance,
	}
	failOnError(commandBus.Handle(context.Background(), openAccountCmd))
}

func transferMoney(commandBus cqrs.CommandHandler, from, to account.Number, amount money.Money) {
	transferCmd := command.TransferMoney{
		From:   from,
		To:     to,
		Amount: amount,
	}
	failOnError(commandBus.Handle(context.Background(), transferCmd))

	receiveCmd := command.ReceiveMoney{
		From:   from,
		To:     to,
		Amount: amount,
	}
	failOnError(commandBus.Handle(context.Background(), receiveCmd))
}

func queryAccount(queryBus cqrs.QueryHandler, number string) {
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
