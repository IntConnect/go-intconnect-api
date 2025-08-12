package exception

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"net/http"
)

func ParseGormError(err error) *ApplicationError {
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return &ApplicationError{
			Message:    "Record not found",
			StatusCode: http.StatusNotFound,
		}

	case errors.Is(err, gorm.ErrDuplicatedKey):
		return &ApplicationError{
			Message:    "Data already exists",
			StatusCode: http.StatusConflict,
		}

	// Handle MySQL/Postgres specific errors
	case errors.Is(err, gorm.ErrForeignKeyViolated):
		return &ApplicationError{
			Message:    "Related record not found",
			StatusCode: http.StatusBadRequest,
		}

	case errors.Is(err, gorm.ErrDuplicatedKey):
		return &ApplicationError{
			Message:    "Duplicate entry",
			StatusCode: http.StatusConflict,
		}
	case errors.Is(err, gorm.ErrInvalidData):
		return &ApplicationError{
			Message:    "Invalid data",
			StatusCode: http.StatusBadRequest,
		}
	case errors.Is(err, gorm.ErrForeignKeyViolated):
		return &ApplicationError{
			Message:    "Terdapat data lain yang berelasi",
			StatusCode: http.StatusBadRequest,
		}
	case errors.Is(err, gorm.ErrCheckConstraintViolated):
		return &ApplicationError{
			Message:    "Check constraint failed",
			StatusCode: http.StatusBadRequest,
		}
	}

	var clientError *ApplicationError
	isApplicationError := errors.As(err, &clientError)
	if isApplicationError {
		return NewApplicationError(clientError.StatusCode, clientError.Message, clientError.rawError)
	}
	return NewApplicationError(http.StatusInternalServerError, "Database error occurred", errors.New("database error occurred"))

}
func ParseGormErrorWithMessage(err error, additionalMessage interface{}) *ApplicationError {
	gormError := ParseGormError(err)
	gormError.Message = fmt.Sprintf("%s %s", gormError.Message, additionalMessage)
	gormError.rawError = err
	gormError.Trace = additionalMessage
	return gormError
}
