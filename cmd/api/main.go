package main

import (
	"github.com/Ayocodes24/GO-Eats/cmd/api/middleware"
	"github.com/Ayocodes24/GO-Eats/pkg/database"
	"github.com/Ayocodes24/GO-Eats/pkg/handler"
	annoucements "github.com/Ayocodes24/GO-Eats/pkg/handler/announcements"
	crt "github.com/Ayocodes24/GO-Eats/pkg/handler/cart"
	delv "github.com/Ayocodes24/GO-Eats/pkg/handler/delivery"
	notify "github.com/Ayocodes24/GO-Eats/pkg/handler/notification"
	"github.com/Ayocodes24/GO-Eats/pkg/handler/restaurant"
	revw "github.com/Ayocodes24/GO-Eats/pkg/handler/review"
	"github.com/Ayocodes24/GO-Eats/pkg/handler/user"
	"github.com/Ayocodes24/GO-Eats/pkg/nats"
	"github.com/Ayocodes24/GO-Eats/pkg/service/announcements"
	"github.com/Ayocodes24/GO-Eats/pkg/service/cart_order"
	"github.com/Ayocodes24/GO-Eats/pkg/service/delivery"
	"github.com/Ayocodes24/GO-Eats/pkg/service/notification"
	restro "github.com/Ayocodes24/GO-Eats/pkg/service/restaurant"
	"github.com/Ayocodes24/GO-Eats/pkg/service/review"
	usr "github.com/Ayocodes24/GO-Eats/pkg/service/user"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
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

	// Connect NATS
	natServer, err := nats.NewNATS("nats://127.0.0.1:4222")

	// WebSocket Clients
	wsClients := make(map[string]*websocket.Conn)

	s := handler.NewServer(db, true)

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

	// Events/Announcements
	announceService := announcements.NewAnnouncementService(db, env)
	annoucements.NewAnnouncementHandler(s, "/announcements", announceService, middlewares, validate)

	// Notification
	notifyService := notification.NewNotificationService(db, env, natServer)

	// Subscribe to multiple events.
	_ = notifyService.SubscribeNewOrders(wsClients)
	_ = notifyService.SubscribeOrderStatus(wsClients)

	notify.NewNotifyHandler(s, "/notify", notifyService, middlewares, validate, wsClients)

	s.Gin.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	log.Fatal(s.Run())

}
