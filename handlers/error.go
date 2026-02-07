package handlers

import "github.com/gin-gonic/gin"

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message`
}

func RespondError(c *gin.Context, status int, code, message string) {
	c.JSON(status, ErrorResponse{
		Error: ErrorDetail{
			Code:    code,
			Message: message,
		},
	})
}
