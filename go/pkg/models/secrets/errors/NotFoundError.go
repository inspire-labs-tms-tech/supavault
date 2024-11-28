package errors

import "fmt"

type NotFoundError struct {
	Err  string
	Hint string
}

func NewNotFoundError() *NotFoundError {
	return &NotFoundError{
		Err:  "credential not found",
		Hint: "are you logged in?",
	}
}

func (e *NotFoundError) Error() string {
	if e.Hint != "" {
		return fmt.Sprintf("%s (%s)", e.Err, e.Hint)
	}
	return e.Err
}
