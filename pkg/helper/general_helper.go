package helper

import (
	"go-intconnect-api/internal/model"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/logger"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func CheckErrorOperation(indicatedError error, applicationError *exception.ApplicationError) bool {
	if indicatedError != nil {
		logger.Debug(indicatedError)
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
	CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrPayloadInvalid))
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
	jwtClaims, isExists := ginContext.Get("claims")
	if !isExists {
		exception.ThrowApplicationError(exception.NewApplicationError(http.StatusUnauthorized, exception.ErrUnauthorized))
	}
	userClaim, isValid := jwtClaims.(*model.JwtClaimRequest)
	if !isValid {
		exception.ThrowApplicationError(exception.NewApplicationError(http.StatusUnauthorized, exception.ErrUnauthorized))
	}

	return userClaim
}

func ExtractRequestMeta(ginContext *gin.Context) (string, string) {
	ipAddress, isIpAddressExists := ginContext.Get("ipAddress")
	userAgent, isUserAgentExists := ginContext.Get("userAgent")
	if isUserAgentExists && isIpAddressExists {
		return ipAddress.(string), userAgent.(string)
	}
	return "", ""
}

func ExtractRequestData(ginContext *gin.Context) (*model.JwtClaimRequest, string, string) {
	userJwtClaims := ExtractJwtClaimFromContext(ginContext)
	ipAddress, userAgent := ExtractRequestMeta(ginContext)
	return userJwtClaims, ipAddress, userAgent
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func TakePointer[T any](value T) *T {
	return &value
}
