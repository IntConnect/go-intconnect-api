package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/helper"
	"go-intconnect-api/pkg/mapper"
	"net/http"
	"strings"
)

func AuthMiddleware(viperConfig *viper.Viper) gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		tokenString := ginContext.GetHeader("Authorization")
		trimmedTokenString := strings.Replace(tokenString, "Bearer ", "", 1)
		// Parse the token
		token, err := jwt.Parse(trimmedTokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrAbortHandler
			}
			return []byte(viperConfig.GetString("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			ginContext.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			ginContext.Abort() // Stop further processing if unauthorized
			return
		}

		// Set the token claims to the context
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userJwtClaim, err := mapper.MapJwtClaimIntoUserClaim(claims)
			helper.CheckErrorOperation(err, exception.NewApplicationError(http.StatusBadRequest, exception.ErrBadRequest, err))
			ginContext.Set("claims", userJwtClaim)
		} else {
			ginContext.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			ginContext.Abort()
			return
		}

		ginContext.Next() // Proceed to the next handler if authorized
	}
}
