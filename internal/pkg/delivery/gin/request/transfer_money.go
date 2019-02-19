package request

type TransferMoney struct {
	To       string `json:"to" example:"ACC555"`
	From     string `json:"from" example:"ACC777"`
	Amount   int64  `json:"amount" example:"10000"`
	Currency string `json:"currency" example:"USD"`
}
