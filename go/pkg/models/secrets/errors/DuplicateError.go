package errors

import "fmt"

type DuplicateError struct {
	Err  string
	Hint string
}

func (e *DuplicateError) Error() string {
	if e.Hint != "" {
		return fmt.Sprintf("%s (%s)", e.Err, e.Hint)
	}
	return e.Err
}
