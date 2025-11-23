package helper

import (
	"go-intconnect-api/pkg/exception"
)

func CheckErrorOperation(indicatedError error, applicationError *exception.ApplicationError) bool {
	if indicatedError != nil {
		panic(applicationError)
		return true
	}
	return false
}

func CheckPointerWrapper[T any](targetChecking *T, renderPayload func()) {
	if targetChecking != nil {
		renderPayload()
	}
}
