package handlers

import (
	"github.com/gin-gonic/gin"
	"shop/internal/services"
)

type Handler struct {
	service *services.Service
}

func NewHandler(s *services.Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) InitRoute() *gin.Engine {
	router := gin.New()

	router.Use(LanguageMiddleware())

	router.GET("/order", h.SelectOrders)
	router.DELETE("/order/:id", h.DeleteOrderById)
	router.POST("/order", h.SaveOrder)

	return router
}
