package exception

import "fmt"

type ApplicationError struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Trace      interface{} `json:"trace"`
	rawError   error
}

func (applicationError *ApplicationError) Error() string {
	return fmt.Sprintf("Error %d: %s", applicationError.StatusCode, applicationError.Message)
}

func NewApplicationError(statusCode int, message string, rawError error) *ApplicationError {
	return &ApplicationError{
		StatusCode: statusCode,
		rawError:   rawError,
		Message:    message,
	}
}

func NewApplicationErrorWithTracing(statusCode int, message string, rawError error, traceError interface{}) *ApplicationError {
	return &ApplicationError{
		StatusCode: statusCode,
		rawError:   rawError,
		Message:    message,
		Trace:      traceError,
	}
}

func (applicationError *ApplicationError) GetRawError() error {
	if applicationError == nil || applicationError.rawError == nil {
		return nil
	}
	return applicationError.rawError
}

func ThrowApplicationError(applicationError *ApplicationError) {
	panic(applicationError)
}
