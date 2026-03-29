package middleware

import (
	"github.com/gin-gonic/gin"
)

var allowedOrigins = map[string]bool{
	"http://127.0.0.1:5500": true,
	"http://localhost:5500": true,

	"https://127.0.0.1:5173": true,
	"https://localhost:5173": true,
}

// CORSMiddleware allows cross-origin requests
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		if allowedOrigins[origin] {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-API-Key")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
