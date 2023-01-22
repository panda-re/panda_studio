package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

var ErrorMapping = map[error]int{
	mongo.ErrNoDocuments: http.StatusNotFound,
}

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


		err := c.Errors[0]
		details := ""

		fmt.Printf("%T %T", err, errors.Cause(err.Err))

		statusCode, ok := ErrorMapping[errors.Cause(err.Err)];
		if !ok {
			statusCode = 500
		}

		if err, ok := err.Err.(stackTracer); ok {
			stack := err.StackTrace()
			details = fmt.Sprintf("%+v", stack)
		}

		errResponse := ErrorMessage{
			Message: fmt.Sprintf("%s", err.Err),
			Details: details,
		}

		// todo: return correct error code
		c.AbortWithStatusJSON(statusCode, gin.H{"error": errResponse })
	}
}