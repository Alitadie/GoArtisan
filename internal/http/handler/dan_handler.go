package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

type danHandler struct {}

func NewdanHandler() *danHandler {
	return &danHandler{}
}

func (h *danHandler) Index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello from dan"})
}
