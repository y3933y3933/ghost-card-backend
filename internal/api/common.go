package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func BadRequest(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, gin.H{"error": msg})
}

func InternalServerError(c *gin.Context, msg string) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
}

func NotFound(c *gin.Context, msg string) {
	c.JSON(http.StatusNotFound, gin.H{"error": msg})

}

func FailedValidation(c *gin.Context, details any) {
	c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Validation failed", "details": details})
}

func Forbidden(c *gin.Context, msg string) {
	c.JSON(http.StatusForbidden, gin.H{"error": msg})

}

func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    data,
	})
}
