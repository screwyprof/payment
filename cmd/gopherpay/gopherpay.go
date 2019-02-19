package main

import (
	"github.com/screwyprof/payment/internal/pkg/app"
)

func main() {
	application := app.NewGopherPay()
	application.Run()
}
