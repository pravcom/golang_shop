package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"shop/internal/models"
)

func (h *Handler) DeleteOrderById(c *gin.Context) {
	id := c.Param("id")

	log.Println("Deleting order with ID: " + id)

	orderId, err := strconv.Atoi(id)

	if err != nil {
		newErrResponse(c, http.StatusBadRequest, "not valid id")
	}

	err = h.service.Order.DeleteById(orderId)

	if err != nil {
		newErrResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Order deleted: " + id,
	})

}

func (h *Handler) SaveOrder(c *gin.Context) {

	var orderRequest models.OrderRequest

	err := c.BindJSON(&orderRequest)
	if err != nil {
		newErrResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	orderResponse, err := h.service.Order.Save(orderRequest)
	if err != nil {
		newErrResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, orderResponse)
}

func (h *Handler) SelectOrders(c *gin.Context) {
	var filter models.OrderFilter

	err := c.ShouldBindQuery(&filter)

	if err != nil {
		newErrResponse(c, http.StatusBadRequest, err.Error())
	}

	lang := c.GetString("lang")
	if lang == "" {
		lang = "ru"
	}

	filter.Lang = &lang

	orders, err := h.service.Select(filter)

	if err != nil {
		newErrResponse(c, http.StatusBadRequest, err.Error())
	}

	c.JSON(http.StatusOK, orders)

}
