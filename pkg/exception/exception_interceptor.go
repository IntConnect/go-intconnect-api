package exception

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-intconnect-api/internal/model"
	"net/http"
)

func Interceptor() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		defer func() {
			if occurredError := recover(); occurredError != nil {
				fmt.Println("panic occurred", occurredError)
				// Check if it's our custom error
				if clientError, ok := occurredError.(*ApplicationError); ok {
					fmt.Println("panic occurred", clientError.GetRawError())
					ginContext.AbortWithStatusJSON(
						clientError.StatusCode,
						model.NewResponseContractModel(false, clientError.Message, nil, &clientError.Trace),
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
