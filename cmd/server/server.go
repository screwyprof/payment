package main

import (
	"github.com/gin-gonic/gin"
	"github.com/screwyprof/payment/pkg/report"

	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"github.com/screwyprof/payment/internal/pkg/delivery/gin/controller"
	_ "github.com/screwyprof/payment/internal/pkg/delivery/gin/docs"

	cmdadaptor "github.com/screwyprof/payment/internal/pkg/adaptor/command_handler"
	qdaptor "github.com/screwyprof/payment/internal/pkg/adaptor/query_handler"

	"github.com/screwyprof/payment/internal/pkg/bus"
	"github.com/screwyprof/payment/internal/pkg/cqrs"
	"github.com/screwyprof/payment/internal/pkg/observer"
	"github.com/screwyprof/payment/internal/pkg/reporting"
	"github.com/screwyprof/payment/internal/pkg/repository"

	"github.com/screwyprof/payment/pkg/command_handler"
	"github.com/screwyprof/payment/pkg/event_handler"
	"github.com/screwyprof/payment/pkg/query_handler"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample Payment service.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url https://github.com/screwyprof/s

// @host localhost:8080
// @BasePath /api/v1

func main() {
	// init deps
	accountReporter := reporting.NewInMemoryAccountReporter()
	accountRepo := repository.NewInMemoryAccountReporter()
	notifier := newNotifier(accountReporter)

	accountOpenner := command_handler.NewOpenAccount(accountRepo, notifier)
	moneyTransfer := command_handler.NewTransferMoney(accountRepo, accountRepo, notifier)
	moneyReceiver := command_handler.NewReceiveMoney(accountRepo, accountRepo, notifier)
	commandBus := newCommandBus(moneyTransfer, moneyReceiver, accountOpenner)

	accountQueryer := query_handler.NewGetAccountShortInfo(accountReporter)
	accountsQueryer := query_handler.NewGetAllAccounts(accountReporter)
	queryBus := newQueryBus(accountQueryer, accountsQueryer)

	openAccountCtrl := controller.NewOpenAccount(commandBus, queryBus)
	showAccountCtrl := controller.NewShowAccount(queryBus)
	transferMoneyCtrl := controller.NewTransferMoney(commandBus, queryBus)
	showAvailableAccountsCtrl := controller.NewShowAvailableAccounts(queryBus)

	// init router
	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		accounts := v1.Group("/accounts")
		{
			accounts.GET(":number", showAccountCtrl.Handle)
			accounts.GET("", showAvailableAccountsCtrl.Handle)
			accounts.POST("", openAccountCtrl.Handle)
			accounts.POST(":number/transfer", transferMoneyCtrl.Handle)
		}
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// run server
	r.Run(":8080")
}

func newNotifier(accountReporter report.AccountUpdater) observer.Notifier {
	notifier := observer.NewNotifier()
	notifier.Register(event_handler.NewAccountOpened(accountReporter))
	notifier.Register(event_handler.NewMoneyTransfered(accountReporter))
	//notifier.Register(event_handler.NewMoneySentFurther(receiveMoney))
	notifier.Register(event_handler.NewMoneyReceived(accountReporter))

	return notifier
}

func newCommandBus(moneyTransfer, moneyReceiver, accountOpenner cqrs.CommandHandler) cqrs.CommandHandler {
	commandBus := bus.NewCommandBus()
	commandBus.Register("OpenAccount", cmdadaptor.FromDomain(accountOpenner))
	commandBus.Register("TransferMoney", cmdadaptor.FromDomain(moneyTransfer))
	commandBus.Register("ReceiveMoney", cmdadaptor.FromDomain(moneyReceiver))

	return cmdadaptor.ToDomain(commandBus)
}

func newQueryBus(accountQueryer cqrs.QueryHandler, accountsQueryer cqrs.QueryHandler) cqrs.QueryHandler {
	queryBus := bus.NewQueryHandlerBus()
	queryBus.Register("GetAccountShortInfo", qdaptor.FromDomain(accountQueryer))
	queryBus.Register("GetAllAccounts", qdaptor.FromDomain(accountsQueryer))

	return qdaptor.ToDomain(queryBus)
}
