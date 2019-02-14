package report

type GetAccountByNumber interface {
	ByNumber(number string) (*Account, error)
}

type AccountUpdater interface {
	Update(account Account) error
}
