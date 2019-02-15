package event_handler

//
//import (
//	"context"
//	"fmt"
//	"github.com/screwyprof/payment/pkg/domain/account"
//
//	"github.com/screwyprof/payment/internal/pkg/cqrs"
//	"github.com/screwyprof/payment/internal/pkg/observer"
//	"github.com/screwyprof/payment/pkg/command"
//	"github.com/screwyprof/payment/pkg/event"
//)
//
//type MoneySentFurther struct {
//	moneyReceiver cqrs.CommandHandler
//}
//
//func NewMoneySentFurther(moneyReceiver cqrs.CommandHandler) *MoneySentFurther {
//	return &MoneySentFurther{moneyReceiver: moneyReceiver}
//}
//
//func (h *MoneySentFurther) Handle(e observer.Event) {
//	evn, ok := e.(event.MoneyTransferred)
//	if !ok {
//		return
//	}
//
//	fmt.Printf("MoneySentFurther: %s=%s, %s => %s\n",
//		evn.From, evn.Balance.Display(), evn.Amount.Display(), evn.To)
//
//	// TODO: if error occurs do a compensating action, i.e MoneyTransferFailedCommand
//	err := h.moneyReceiver.Handle(context.Background(), command.ReceiveMoney{
//		From:   account.Number(evn.From),
//		To:     account.Number(evn.To),
//		Amount: evn.Amount,
//	})
//
//	if err != nil {
//		panic(err)
//	}
//}
