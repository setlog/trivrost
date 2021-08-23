package system

type universalNameRetrievalErrorType int

const errorBadDevice universalNameRetrievalErrorType = 1
const errorConnectionUnavailable universalNameRetrievalErrorType = 2
const errorExtendedError universalNameRetrievalErrorType = 3
const errorMoreData universalNameRetrievalErrorType = 4
const errorNotSupported universalNameRetrievalErrorType = 5
const errorNoNetOrBadPath universalNameRetrievalErrorType = 6
const errorNoNetwork universalNameRetrievalErrorType = 7
const errorNotConnected universalNameRetrievalErrorType = 8
const errorUndocumented universalNameRetrievalErrorType = 9

type universalNameRetrievalError struct {
	message   string
	errorType universalNameRetrievalErrorType
}

func (err *universalNameRetrievalError) Error() string {
	if err == nil {
		return "<nil>"
	}
	return err.message
}

// ErrorType returns the corresponsing WINAPI error type of the WNetGetUniversalNameW function call which generated the error.
func (err *universalNameRetrievalError) ErrorType() universalNameRetrievalErrorType {
	return err.errorType
}
