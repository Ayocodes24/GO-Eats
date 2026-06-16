package restaurant

import (
	"context"
	"github.com/Ayocodes24/GO-Eats/pkg/database/models/restaurant"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func (s *RestaurantHandler) addMenu(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	var menuItem restaurant.MenuItem
	if err := c.BindJSON(&menuItem); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	photoAlreadySet := menuItem.Photo != ""

	menuObject, err := s.service.AddMenu(ctx, &menuItem)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Only call Unsplash if no photo was provided in the request
	if !photoAlreadySet {
		s.service.UpdateMenuPhoto(ctx, menuObject)
	}

	c.JSON(http.StatusCreated, gin.H{"message": "New Menu Added!"})
}

func (s *RestaurantHandler) listMenus(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	restaurantId := c.Query("restaurant_id")
	if restaurantId == "" {
		results, err := s.service.ListAllMenus(ctx)
		if err != nil {
			c.JSON(http.StatusOK, []restaurant.MenuItem{})
			return
		}
		c.JSON(http.StatusOK, results)
		return
	}

	restaurantID, err := strconv.ParseInt(restaurantId, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid restaurant_id"})
		return
	}

	results, err := s.service.ListMenus(ctx, restaurantID)
	if err != nil || len(results) == 0 {
		c.JSON(http.StatusOK, []restaurant.MenuItem{})
		return
	}
	c.JSON(http.StatusOK, results)
}

func (s *RestaurantHandler) deleteMenu(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	menuId, err := strconv.ParseInt(c.Param("menu_id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid menu_id"})
		return
	}
	restaurantId, err := strconv.ParseInt(c.Param("restaurant_id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid restaurant_id"})
		return
	}

	if _, err = s.service.DeleteMenu(ctx, menuId, restaurantId); err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
