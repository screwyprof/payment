package report

import "github.com/google/uuid"

type GetAccountByNumber interface {
	ByNumber(number string) (*Account, error)
}

type IDByNumber interface {
	IDByNumber(number string) (uuid.UUID, error)
}

type GetAllAccounts interface {
	All() (Accounts, error)
}

type AccountUpdater interface {
	Update(account *Account) error
}
