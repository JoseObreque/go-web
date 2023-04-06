package middleware

import (
	"errors"
	"github.com/JoseObreque/go-web/pkg/web"
	"github.com/gin-gonic/gin"
	"os"
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
