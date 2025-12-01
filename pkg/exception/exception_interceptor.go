package exception

import (
	"go-intconnect-api/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Interceptor() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		defer func() {
			if occurredError := recover(); occurredError != nil {
				logrus.Debug("panic occurred", occurredError)
				// Check if it's our custom error
				if clientError, ok := occurredError.(*ApplicationError); ok {
					ginContext.AbortWithStatusJSON(
						clientError.HttpStatusCode,
						&model.ResponseContract[any]{
							Success: false,
							Message: "",
							Entry:   nil,
							Entries: nil,
							Error: &model.ErrorDetail{
								Code:    clientError.ConventionStatusCode,
								Message: clientError.Message,
								Details: clientError.Details,
							},
						},
					)
					return
				}

				// Unknown error
				ginContext.AbortWithStatusJSON(
					http.StatusInternalServerError,
					model.NewResponseContractModel(false, "Internal server error", nil, nil),
				)
			}
		}()
		ginContext.Next()
	}
}
