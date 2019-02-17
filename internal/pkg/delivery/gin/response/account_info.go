package response

type Ledger struct {
	Action string `json:"action" example:"Transfer from AK777, $100"`
}

type AccountInfo struct {
	Number  string `json:"number" example:"ACC777"`
	Balance string `json:"balance" example:"$100.00"`

	Ledgers []Ledger `json:"ledgers"`
}
