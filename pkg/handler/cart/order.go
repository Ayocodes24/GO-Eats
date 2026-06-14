package cart

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func (s *CartHandler) PlaceNewOrder(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var req struct {
		DeliveryAddress string `json:"delivery_address" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "delivery_address is required"})
		return
	}

	userID := c.GetInt64("userID")
	cartInfo, err := s.service.GetCartId(ctx, userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newOrder, err := s.service.PlaceOrder(ctx, cartInfo.CartID, userID, req.DeliveryAddress)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := s.service.RemoveItemsFromCart(ctx, cartInfo.CartID); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := s.service.NewOrderPlacedNotification(userID, newOrder.OrderID); err != nil {
		slog.Warn("order notification failed", "order_id", newOrder.OrderID, "error", err)
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Order placed!"})
}

func (s *CartHandler) getOrderList(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userID := c.GetInt64("userID")
	orders, err := s.service.OrderList(ctx, userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"orders": orders})
	return
}

func (s *CartHandler) getOrderItemsList(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userID := c.GetInt64("userID")
	orderID, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	orders, err := s.service.OrderItemsList(ctx, userID, orderID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"orders": orders})
	return
}

func (s *CartHandler) getDeliveriesList(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userID := c.GetInt64("userID")
	orderID, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	deliveries, err := s.service.DeliveryInformation(ctx, orderID, userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"delivery_info": deliveries})
	return
}
