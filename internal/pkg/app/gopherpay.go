package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"github.com/screwyprof/payment/internal/pkg/delivery/gin/controller"
	_ "github.com/screwyprof/payment/internal/pkg/delivery/gin/docs"

	qdaptor "github.com/screwyprof/payment/internal/pkg/adaptor/query_handler"

	"github.com/screwyprof/payment/internal/pkg/bus"
	"github.com/screwyprof/payment/internal/pkg/cqrs"
	"github.com/screwyprof/payment/internal/pkg/observer"
	"github.com/screwyprof/payment/internal/pkg/reporting"

	"github.com/screwyprof/payment/pkg/domain"
	"github.com/screwyprof/payment/pkg/domain/account"
	"github.com/screwyprof/payment/pkg/event_handler"
	"github.com/screwyprof/payment/pkg/query_handler"
	"github.com/screwyprof/payment/pkg/report"
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

type GopherPay struct {
	srv    *http.Server
	router *gin.Engine
}

func NewGopherPay() *GopherPay {
	router := SetupRouter()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	return &GopherPay{router: router, srv: srv}
}

func (g *GopherPay) Run() {
	go func() {
		if err := g.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	<-quit
	g.Shutdown()
}

func (g *GopherPay) Shutdown() {
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := g.srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}

func SetupRouter() *gin.Engine {
	// init deps
	cqrs.RegisterAggregate(func(id uuid.UUID) domain.Aggregate {
		return account.Create(id)
	})

	accountReporter := reporting.NewInMemoryAccountReporter()
	notifier := newNotifier(accountReporter)

	eventStore := cqrs.NewInMemoryEventStore()
	commandHandler := cqrs.NewEventSourceHandler(eventStore, notifier)

	accountQueryer := query_handler.NewGetAccountShortInfo(accountReporter)
	availableAccountsQueryer := query_handler.NewGetAllAccounts(accountReporter)
	queryBus := newQueryBus(accountQueryer, availableAccountsQueryer)

	openAccountCtrl := controller.NewOpenAccount(commandHandler, queryBus)
	showAccountCtrl := controller.NewShowAccount(queryBus)
	transferMoneyCtrl := controller.NewTransferMoney(commandHandler, accountReporter)
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

	return r
}

func newNotifier(accountReporter report.AccountUpdater) observer.Notifier {
	notifier := observer.NewNotifier()
	notifier.Register(event_handler.NewAccountOpened(accountReporter))
	notifier.Register(event_handler.NewMoneyTransfered(accountReporter))
	notifier.Register(event_handler.NewMoneyReceived(accountReporter))

	return notifier
}

func newQueryBus(accountQueryer, availableAccountsQueryer domain.QueryHandler) domain.QueryHandler {
	queryBus := bus.NewQueryHandlerBus()
	queryBus.Register("GetAccountShortInfo", qdaptor.FromDomain(accountQueryer))
	queryBus.Register("GetAllAccounts", qdaptor.FromDomain(availableAccountsQueryer))

	return qdaptor.ToDomain(queryBus)
}
