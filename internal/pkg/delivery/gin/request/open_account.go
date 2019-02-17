package request

type OpenAccount struct {
	Number   string `json:"number" example:"ACC777"`
	Amount   int64  `json:"amount" example:"77700"`
	Currency string `json:"currency" example:"USD"`
}

//func (a OpenAccount) Validate() error {
// Check that number has valid format
// Check that Amount is Positive and non Zero
// Check that Currency is correct.
//}
