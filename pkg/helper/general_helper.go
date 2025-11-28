package helper

import (
	"go-intconnect-api/internal/model"
	"go-intconnect-api/pkg/exception"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
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

func ExtractIndexedFiles(
	ginContext *gin.Context,
	prefixKey string,
	suffixKey string,
	expectedLen int,
) ([]*multipart.FileHeader, error) {

	multipartForm, err := ginContext.MultipartForm()
	CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))
	fileHeadersArray := make([]*multipart.FileHeader, expectedLen)

	for key, fileHeaders := range multipartForm.File {
		if strings.HasPrefix(key, prefixKey) &&
			strings.HasSuffix(key, suffixKey) {

			idxStr := key[len(prefixKey):strings.Index(key, "]")]
			idxInt, err := strconv.Atoi(idxStr)
			if err != nil || idxInt < 0 || idxInt >= expectedLen {
				continue
			}

			if len(fileHeaders) > 0 {
				fileHeadersArray[idxInt] = fileHeaders[0]
			}
		}
	}

	return fileHeadersArray, nil
}

func ExtractJwtClaimFromContext(ginContext *gin.Context) *model.JwtClaimRequest {
	claims, exists := ginContext.Get("claims")
	if !exists {
		exception.ThrowApplicationError(exception.NewApplicationError(http.StatusUnauthorized, exception.ErrUnauthorized, nil))
	}
	userClaim, ok := claims.(*model.JwtClaimRequest)
	if !ok {
		exception.ThrowApplicationError(exception.NewApplicationError(http.StatusUnauthorized, exception.ErrUnauthorized, nil))
	}
	return userClaim
}
