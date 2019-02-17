package response

type ShortAccountInfo struct {
	Number  string `json:"number" example:"ACC777"`
	Balance string `json:"balance" example:"$100.00"`
}
