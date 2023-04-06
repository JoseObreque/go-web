package middleware

import (
	"errors"
	"github.com/JoseObreque/go-web/pkg/web"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"time"
)

var ErrInvalidToken = errors.New("invalid token")

func TokenValidator() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the token from the request header
		token := c.GetHeader("token")

		// Check if the token is not empty
		if token == "" {
			c.Abort()
			web.Failure(c, 401, ErrInvalidToken)
			return
		}

		// Check if the token is valid
		if token != os.Getenv("TOKEN") {
			c.Abort()
			web.Failure(c, 401, ErrInvalidToken)
			return
		}

		c.Next()
	}
}

func PanicLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				now := time.Now()
				log.Printf("HTTP Verb: %s\n", c.Request.Method)
				log.Printf("URL: %s\n", c.Request.URL.Path)
				log.Printf("Datetime: %s\n", now.Format("2006-01-02 15:04:05"))
				log.Printf("Bytes: %b\n", c.Request.ContentLength)
			}
		}()

		c.Next()
	}
}
