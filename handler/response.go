package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

func sendErrorResponse(c *gin.Context, statusCode int, message, errorType string) {
	c.JSON(statusCode, ErrorResponse{
		Message: message,
		Error:   errorType,
	})
}

func BadRequest(c *gin.Context, message string) {
	sendErrorResponse(c, http.StatusBadRequest, message, "Bad Request")
}

func NotFound(c *gin.Context, message string) {
	sendErrorResponse(c, http.StatusNotFound, message, "Not Found")
}

func InternalServerError(c *gin.Context, message string) {
	sendErrorResponse(c, http.StatusInternalServerError, message, "Internal Server Error")
}

func OK(c *gin.Context, message string) {
	c.JSON(http.StatusOK, gin.H{"message": message})
}
