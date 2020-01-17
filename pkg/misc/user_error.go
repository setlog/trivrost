package misc

import (
	"fmt"
)

const programmingErrorInfo = "A programming error has been encountered."
const uncommonErrorPrefix = "An error occurred. Technical information: "
const unknownErrorInfo = "An unknown error occurred."

// UserError is an error with a cause-agnostic message and the added semantic of that message being suitable
// for display to users, i.e. non-programmers. It allows you to construct an error which will return a cause
// with Unwrap(), without the unwrapped error's message appearing in strings returned from Message().
type UserError struct {
	cause       error
	userMessage string
}

func (e *UserError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.cause == nil {
		return e.Message()
	}
	return e.Message() + " Cause: " + e.cause.Error()
}

func (e *UserError) Unwrap() error {
	return e.cause
}

// Message returns a user-readable explanation of the error.
func (e *UserError) Message() string {
	if e == nil {
		return programmingErrorInfo
	}
	if e.userMessage == "" {
		if e.cause != nil {
			return uncommonErrorPrefix + e.cause.Error()
		}
		return unknownErrorInfo
	}
	return e.userMessage
}

// UserErrorf constructs a new *UserError with given cause and formatted message. The message should start with
// a capital letter and have proper punctuation - but no word-wrapping line-breaks - for display in a dialog.
func UserErrorf(cause error, userMessageFormat string, args ...interface{}) error {
	return &UserError{cause: cause, userMessage: fmt.Sprintf(userMessageFormat, args...)}
}

// NewUserErrorFromErrors returns the first non-nil error in the variadic parameter list if it is a *UserError. If the first
// non-nil error is not a *UserError, one is built using UserErrorf() with a generic message revealing the result of Error().
func NewUserErrorFromErrors(errs ...error) error {
	for _, err := range errs {
		if !IsNil(err) {
			if userError, isUserError := err.(*UserError); isUserError {
				return userError
			}
			return UserErrorf(err, "%s%v", uncommonErrorPrefix, err)
		}
	}
	return nil
}
