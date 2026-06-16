package restaurant

import (
	"context"
	restaurantModel "github.com/Ayocodes24/GO-Eats/pkg/database/models/restaurant"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func (s *RestaurantHandler) addRestaurant(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	_ = ctx

	var photoPath string

	// File is optional — if provided, upload it; otherwise use image_url form field
	file, fileHeader, err := c.Request.FormFile("file")
	if err == nil {
		newFileName := generateFileName(fileHeader.Filename)
		if _, uploadErr := s.Serve.Storage.Upload(newFileName, file); uploadErr != nil {
			slog.Error("addRestaurant: upload failed", "error", uploadErr)
		}
		photoPath = filepath.Join(os.Getenv("STORAGE_DIRECTORY"), newFileName)
	} else {
		photoPath = c.PostForm("image_url")
	}

	var restaurant restaurantModel.Restaurant
	restaurant.Name        = c.PostForm("name")
	restaurant.Description = c.PostForm("description")
	restaurant.Address     = c.PostForm("address")
	restaurant.City        = c.PostForm("city")
	restaurant.State       = c.PostForm("state")
	restaurant.Photo       = photoPath

	if _, err := s.service.Add(ctx, &restaurant); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Restaurant created successfully"})
}

func (s *RestaurantHandler) listRestaurants(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	results, err := s.service.ListRestaurants(ctx)
	if err != nil || results == nil {
		// Return empty array instead of 404 when no restaurants exist
		c.JSON(http.StatusOK, []restaurantModel.Restaurant{})
		return
	}
	c.JSON(http.StatusOK, results)
}

func (s *RestaurantHandler) listRestaurantById(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	restaurantId := c.Param("id")
	restaurantID, _ := strconv.ParseInt(restaurantId, 10, 64)

	result, err := s.service.ListRestaurantById(ctx, restaurantID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (s *RestaurantHandler) deleteRestaurant(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	restaurantId  := c.Param("id")
	restaurantID, _ := strconv.ParseInt(restaurantId, 10, 64)

	if _, err := s.service.DeleteRestaurant(ctx, restaurantID); err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
