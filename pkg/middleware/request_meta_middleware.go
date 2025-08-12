package middleware

import "github.com/gin-gonic/gin"

func RequestMetaMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		userAgent := c.GetHeader("User-Agent")

		// Simpan ke context
		c.Set("ipAddress", ip)
		c.Set("userAgent", userAgent)

		c.Next()
	}
}
