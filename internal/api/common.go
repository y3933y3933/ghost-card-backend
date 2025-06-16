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

func NotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{"error": "the requested resource could not be found"})

}

func FailedValidation(c *gin.Context, details any) {
	c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Validation failed", "details": details})
}
