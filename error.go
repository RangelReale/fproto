package fproto

// This error is issued when the protobuf file is malformed
type InvalidScope struct {
	message string
}

func (e *InvalidScope) Error() string {
	return e.message
}
