package query

type GetAccountShortInfo struct {
	Number string
}

func (r GetAccountShortInfo) QueryID() string {
	return "GetAccountShortInfo"
}
