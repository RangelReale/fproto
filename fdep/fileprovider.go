package fdep

import "io"

type FileProvider interface {
	// Return if have more files
	HasNext() bool

	// Returns the current file, and advance the internal pointer
	GetNext() (filepath string, r io.Reader, deptype FileDepType, err error)
}
