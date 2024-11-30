package supabase

import "fmt"

type ClientError struct {
	Msg string
	Err error
}

func NewClientError(msg string, err error) *ClientError {
	return &ClientError{
		Msg: msg,
		Err: err,
	}
}

func (e *ClientError) Error() string {
	if e.Msg != "" {
		return fmt.Sprintf("%s (%s)", e.Msg, e.Err.Error())
	}
	return e.Err.Error()
}
