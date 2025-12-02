package middleware

import (
	"context"
	"fmt"
	"go-intconnect-api/configs"
	"go-intconnect-api/internal/model"
	"go-intconnect-api/pkg/exception"
	"go-intconnect-api/pkg/helper"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func AuthMiddleware(viperConfig *viper.Viper, redisConfig *configs.RedisInstance) gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		// 1. Ambil token dari header
		authHeader := ginContext.GetHeader("Authorization")
		if authHeader == "" {
			ginContext.JSON(http.StatusUnauthorized, helper.NewErrorResponse("", helper.NewErrorDetail(
				exception.StatusAuthError, exception.MsgAuthError, nil)))
			ginContext.Abort()
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		// 2. Parse token untuk ambil userId
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrAbortHandler
			}
			return []byte(viperConfig.GetString("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			ginContext.JSON(http.StatusUnauthorized, helper.NewErrorResponse("", helper.NewErrorDetail(
				exception.StatusAuthError, exception.MsgAuthError, nil)))
			ginContext.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			ginContext.JSON(http.StatusUnauthorized, helper.NewErrorResponse("", helper.NewErrorDetail(
				exception.StatusAuthError, exception.MsgAuthError, nil)))
			ginContext.Abort()
			return
		}

		// 3. Ambil token valid dari Redis
		userId := int64(claims["id"].(float64)) // asumsi "id" selalu ada
		redisKey := fmt.Sprintf("auth:token:%d", userId)
		cachedToken, err := redisConfig.RedisClient.Get(context.Background(), redisKey).Result()
		fmt.Println(err, cachedToken, redisKey, userId)
		if err == redis.Nil || cachedToken != tokenString {
			ginContext.JSON(http.StatusUnauthorized, helper.NewErrorResponse("", helper.NewErrorDetail(
				exception.StatusAuthError, exception.ErrTokenExpiredOrInvalid, nil)))
			ginContext.Abort()
			return
		} else if err != nil {
			ginContext.JSON(http.StatusInternalServerError, helper.NewErrorResponse("", helper.NewErrorDetail(
				exception.StatusInternalError, "Redis error/down", nil)))

			ginContext.Abort()
			return
		}

		// 4. Set claims di context
		userJwtClaim := helper.MapCreateRequestIntoEntity[jwt.MapClaims, model.JwtClaimRequest](&claims)
		ginContext.Set("claims", userJwtClaim)

		ginContext.Next()
	}
}
