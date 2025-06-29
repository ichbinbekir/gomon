package gomon

import "fmt"

type SaveError struct {
	Op  string // truncate, seek, encode, vb.
	Err error
}

func (e *SaveError) Error() string {
	return fmt.Sprintf("save error during %s: %v", e.Op, e.Err)
}

func (e *SaveError) Unwrap() error {
	return e.Err
}
