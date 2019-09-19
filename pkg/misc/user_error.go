package misc

type IUserError interface {
	error
	UserError() string
}

type NestedError struct {
	message string
	cause   error
}

func (e *NestedError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.cause == nil {
		return e.message
	}
	return e.message + ": " + e.cause.Error()
}

func (e *NestedError) Cause() error {
	if e == nil {
		return nil
	}
	return e.cause
}

func (e *NestedError) UserError() string {
	if e == nil {
		return "An unknown error occurred."
	}
	return e.message
}

func NewNestedError(message string, cause error) error {
	return &NestedError{message: message, cause: cause}
}
