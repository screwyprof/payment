package report

type GetAccountByNumber interface {
	ByNumber(number string) (*Account, error)
}

type GetAllAccounts interface {
	All() (Accounts, error)
}

type AccountUpdater interface {
	Update(account *Account) error
}
