package query

type GetAllAccounts struct{}

func (r GetAllAccounts) QueryID() string {
	return "GetAllAccounts"
}
