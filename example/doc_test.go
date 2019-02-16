package example_test

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

func Example() {
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

	// Output:
	// WriteSide: saving account ACC500 = $500.00
	// WriteSide: saving account ACC300 = $300.00
	// ReadSide: retreiving account ACC500
	// ReadSide: account ACC500 retrieved with balance $500.00
	// #ACC500: $500.00
	// Ledgers:
	// #1. Deposit, $500.00
	// ReadSide: retreiving account ACC300
	// ReadSide: account ACC300 retrieved with balance $300.00
	// #ACC300: $300.00
	// Ledgers:
	// #1. Deposit, $300.00
	// WriteSide: retrieving account ACC500
	// WriteSide: account ACC500 retrieved with balance $500.00
	// WriteSide: saving account ACC500 = $200.00
	// ReadSide: updating account ACC500 with balance $200.00
	// WriteSide: retrieving account ACC300
	// WriteSide: account ACC300 retrieved with balance $300.00
	// WriteSide: saving account ACC300 = $600.00
	// ReadSide: updating account ACC300 with balance $600.00
	// ReadSide: retreiving account ACC500
	// ReadSide: account ACC500 retrieved with balance $200.00
	// #ACC500: $200.00
	// Ledgers:
	// #1. Deposit, $500.00
	// #2. Transfer to ACC300, $300.00
	// ReadSide: retreiving account ACC300
	// ReadSide: account ACC300 retrieved with balance $600.00
	// #ACC300: $600.00
	// Ledgers:
	// #1. Deposit, $300.00
	// #2. Transfer from ACC500, $300.00

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
