package messaging

import "fmt"

type ErrConnect struct {
	msg string
	err error
}

func NewErrConnect(err error) *ErrConnect {
	return &ErrConnect{
		msg: "got error during connection",
		err: err,
	}
}

func (e *ErrConnect) Error() string {
	return fmt.Sprintf("%s: %v", e.msg, e.err)
}

func (e *ErrConnect) Unwrap() error {
	return e.err
}

type ErrReceive struct {
	msg string
	err error
}

func NewErrReceive(err error) *ErrReceive {
	return &ErrReceive{
		msg: "got error while sending message",
		err: err,
	}
}

func (e *ErrReceive) Error() string {
	return fmt.Sprintf("%s: %v", e.msg, e.err)
}

func (e *ErrReceive) Unwrap() error {
	return e.err
}

type ErrSend struct {
	msg string
	err error
}

func NewErrSend(err error) *ErrSend {
	return &ErrSend{
		msg: "got error while receiving message",
		err: err,
	}
}

func (e *ErrSend) Error() string {
	return fmt.Sprintf("%s: %v", e.msg, e.err)
}

func (e *ErrSend) Unwrap() error {
	return e.err
}
