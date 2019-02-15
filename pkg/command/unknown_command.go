package command

type Unknown struct{}

func (r Unknown) CommandID() string {
	return "Unknown"
}
