package misc_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/setlog/trivrost/pkg/misc"
)

func TestUserError(t *testing.T) {
	cause := fmt.Errorf("A cause")
	problemFormat := "A problem: %s."
	problemArg := "very technical problem"
	problem := fmt.Sprintf(problemFormat, problemArg)
	e := misc.UserErrorf(cause, problemFormat, problemArg)
	eError := e.Error()
	eErrorExpected := fmt.Sprintf("%s Cause: %s", problem, cause)
	if eError != eErrorExpected {
		t.Errorf("Error() returned \"%s\". Expected: \"%s\"", eError, eErrorExpected)
	}
	eCause := errors.Unwrap(e)
	if eCause != cause {
		t.Errorf("errors.Unwrap() returned \"%v\". Expected: \"%v\"", eCause, cause)
	}
	eMessage := e.(*misc.UserError).Message()
	if eMessage != problem {
		t.Errorf("Message() returned \"%v\". Expected: \"%v\"", eMessage, problem)
	}
	e = nil
	eCause = errors.Unwrap(e)
	if eCause != nil {
		t.Errorf("errors.Unwrap() returned non-nil error for nilled variable.")
	}
}
