package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

var ErrorMapping = map[error]int{
	// MongoDB could not find documents
	mongo.ErrNoDocuments: http.StatusNotFound,
	// Trying to upload a file
	http.ErrNotMultipart: http.StatusBadRequest,
}

func getErrorCode(err error) int {
	code, ok := ErrorMapping[err]
	if !ok {
		code = 500
	}
	return code
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
		if ok, panicVal := func() (ok bool, panicVal any) {
			defer func() {
				if r := recover(); r != nil {
					panicVal = r
					ok = false
				}
			}()

			c.Next()

			return true, nil
		}(); !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": ErrorMessage{
					Message: fmt.Sprintf("%+v", panicVal),
				},
			})
		}

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors[0]
		details := ""

		statusCode := 500//getErrorCode(errors.Cause(err.Err))

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