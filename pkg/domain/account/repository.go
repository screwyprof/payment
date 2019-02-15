package account

type GetAccountByNumber interface {
	ByNumber(number Number) (*Account, error)
}

type AccountStorage interface {
	Add(acc *Account) error
}

//type EventStore interface {
//	StoreEvent(e cqrs.Event)
//}
