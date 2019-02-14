package query

type Unknown struct{}

func (r Unknown) QueryID() string {
	return "Unknown"
}
