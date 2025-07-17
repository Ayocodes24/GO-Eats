package main

import (
	"github.com/Ayocodes24/GO-Eats/cmd/api/middleware"
	"github.com/Ayocodes24/GO-Eats/pkg/abstract/config"
	"github.com/Ayocodes24/GO-Eats/pkg/database"
	crt "github.com/Ayocodes24/GO-Eats/pkg/handler/cart"
	"github.com/Ayocodes24/GO-Eats/pkg/handler/restaurant"
	revw "github.com/Ayocodes24/GO-Eats/pkg/handler/review"
	"github.com/Ayocodes24/GO-Eats/pkg/handler/user"
	"github.com/Ayocodes24/GO-Eats/pkg/service/cart_order"
	restro "github.com/Ayocodes24/GO-Eats/pkg/service/restaurant"
	"github.com/Ayocodes24/GO-Eats/pkg/service/review"
	usr "github.com/Ayocodes24/GO-Eats/pkg/service/user"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"golang.org/x/net/websocket"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	env := os.Getenv("APP_ENV")
	db := database.New()
	// Create Tables
	if err := db.Migrate(); err != nil {
		log.Fatalf("Error migrating database: %s", err)
	}

	// Initialize Validator
	validate := validator.New()

	// Middlewares List
	middlewares := []gin.HandlerFunc{middleware.AuthMiddleware()}

	// User
	userService := usr.NewUserService(db, env)
	user.NewUserHandler(s, "/user", userService, validate)

	// Reviews
	reviewService := review.NewReviewService(db, env)
	revw.NewReviewProtectedHandler(s, "/review", reviewService, middlewares, validate)

	// Restaurant
	restaurantService := restro.NewRestaurantService(db, env)
	restaurant.NewRestaurantHandler(s, "/restaurant", restaurantService)

	// Cart
	cartService := cart_order.NewCartService(db, env, natServer)
	crt.NewCartHandler(s, "/cart", cartService, middlewares, validate)

	// Delivery
	deliveryService := delivery.NewDeliveryService(db, env, natServer)
	delv.NewDeliveryHandler(s, "/delivery", deliveryService, middlewares, validate)

}
