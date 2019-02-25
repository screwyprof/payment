package command

type Unknown struct{}

func (r Unknown) CommandID() string {
	return "Unknown"
}

func (c Unknown) AggregateType() string {
	return "unknown.Unknown"
}
