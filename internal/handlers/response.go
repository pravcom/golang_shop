package handlers

import (
	"log"

	"github.com/gin-gonic/gin"
)

type errResponse struct {
	Message string `json:"message"`
}

type statusResponse struct {
	Status string `json:"status"`
}

func newErrResponse(c *gin.Context, statusCode int, message string) {
	log.Println(message)
	c.AbortWithStatusJSON(statusCode, errResponse{message})
}
