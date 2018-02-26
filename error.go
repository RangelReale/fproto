package fproto

type InvalidScope struct {
	message string
}

func (e *InvalidScope) Error() string {
	return e.message
}
