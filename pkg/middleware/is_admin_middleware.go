package middleware

import (
	"go-intconnect-api/internal/model"
	"go-intconnect-api/pkg/exception"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// RBACMiddleware mengecek apakah user memiliki permission yang dibutuhkan
func IsAdminMiddleware() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		// Ambil claims dari context
		claims, exists := ginContext.Get("claims")
		if !exists {
			ginContext.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			ginContext.Abort()
			return
		}

		// Konversi claims ke struct userJwtClaim
		userClaim, ok := claims.(*model.JwtClaimRequest)
		if !ok {
			exception.ThrowApplicationError(exception.NewApplicationError(http.StatusForbidden, "Invalid token claim"))
			return
		}

		// Cek apakah user memiliki permission yang sesuai
		if userClaim.Role != "Admin" {
			exception.ThrowApplicationError(exception.NewApplicationError(http.StatusForbidden, "Forbidden"))
			return
		}

		// Lanjut ke handler jika punya permission
		ginContext.Next()
	}
}

// hasPermission mengecek apakah permission yang dibutuhkan ada dalam daftar user
func hasAnyPermission(userPermissions []string, requiredPermissions []string) bool {
	for _, required := range requiredPermissions {
		for _, userPerm := range userPermissions {
			if strings.EqualFold(userPerm, required) { // Case-insensitive comparison
				return true
			}
		}
	}
	return false
}
