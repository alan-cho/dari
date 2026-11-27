package errors

import (
	stderrors "errors"
	"fmt"
)

type Error struct {
	Code  Code
	Op    string
	Cause error
}

func New(code Code, op string, cause error) error {
	return &Error{Code: code, Op: op, Cause: cause}
}

func (e *Error) Description() string {
	return describe(e.Code)
}

func (e *Error) Error() string {
	message := fmt.Sprintf("%s: [%s] %s", e.Op, e.Code, e.Description())
	if e.Cause == nil {
		return message
	}
	return fmt.Sprintf("%s: %v", message, e.Cause)
}

func (e *Error) Unwrap() error {
	return e.Cause
}

func HasCode(err error, code Code) bool {
	var harnessErr *Error
	return stderrors.As(err, &harnessErr) && harnessErr.Code == code
}
