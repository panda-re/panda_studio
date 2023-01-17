package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type ErrorMessage struct {
	Message string 	`json:"message"`
	Details string	`json:"stack,omitempty"`
}



func ErrorHandler() gin.HandlerFunc {
	type stackTracer interface {
		StackTrace() errors.StackTrace
	}

	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}
		
		errs := make([]ErrorMessage, len(c.Errors))
		for i, err := range c.Errors {
			var details string = ""

			if err, ok := err.Err.(stackTracer); ok {
				stack := err.StackTrace()
				details = fmt.Sprintf("%+v", stack)
			}


			errs[i] = ErrorMessage{
				Message: fmt.Sprintf("%s", err.Err),
				Details: details,
			}
		}

		// todo: return correct error code
		c.AbortWithStatusJSON(500, gin.H{"errors": errs })
	}
}