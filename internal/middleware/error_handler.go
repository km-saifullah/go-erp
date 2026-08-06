package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	appErrors "github.com/km-saifullah/go-erp/internal/errors"
	"github.com/km-saifullah/go-erp/internal/response"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		var appErr *appErrors.AppError

		if errors.As(err, &appErr) {
			c.JSON(appErr.HTTPStatus, response.ErrorResponse{
				Success: false,
				Error: response.Error{
					Code:    appErr.Code,
					Message: appErr.Message,
					Details: appErr.Details,
				},
			})

			return
		}

		// Unknown/unhandled error.
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Success: false,
			Error: response.Error{
				Code:    appErrors.CodeInternalServer,
				Message: "An unexpected error occurred",
			},
		})
	}
}
