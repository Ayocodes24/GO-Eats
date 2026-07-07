package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/Ayocodes24/GO-Eats/pkg/database"
	cartmodel "github.com/Ayocodes24/GO-Eats/pkg/database/models/cart"
	deliverymodel "github.com/Ayocodes24/GO-Eats/pkg/database/models/delivery"
	ordermodel "github.com/Ayocodes24/GO-Eats/pkg/database/models/order"
	reviewmodel "github.com/Ayocodes24/GO-Eats/pkg/database/models/review"
	"github.com/Ayocodes24/GO-Eats/pkg/handler"
	annhandler "github.com/Ayocodes24/GO-Eats/pkg/handler/announcements"
	carthandler "github.com/Ayocodes24/GO-Eats/pkg/handler/cart"
	deliveryhandler "github.com/Ayocodes24/GO-Eats/pkg/handler/delivery"
	notifyhandler "github.com/Ayocodes24/GO-Eats/pkg/handler/notification"
	reviewhandler "github.com/Ayocodes24/GO-Eats/pkg/handler/review"
	"github.com/Ayocodes24/GO-Eats/pkg/nats"
	restaurantpb "github.com/Ayocodes24/GO-Eats/pkg/proto/restaurant"
	userpb "github.com/Ayocodes24/GO-Eats/pkg/proto/user"
	annsvc "github.com/Ayocodes24/GO-Eats/pkg/service/announcements"
	cartsvc "github.com/Ayocodes24/GO-Eats/pkg/service/cart_order"
	deliverysvc "github.com/Ayocodes24/GO-Eats/pkg/service/delivery"
	notifysvc "github.com/Ayocodes24/GO-Eats/pkg/service/notification"
	reviewsvc "github.com/Ayocodes24/GO-Eats/pkg/service/review"
	"github.com/Ayocodes24/GO-Eats/pkg/wsclients"
	ordermiddleware "github.com/Ayocodes24/GO-Eats/cmd/order-svc/middleware"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	if err := godotenv.Load(".env.order"); err != nil {
		log.Fatal("Error loading .env.order")
	}

	db := database.New()
	if err := db.MigrateModels([]interface{}{
		(*reviewmodel.Review)(nil),
		(*ordermodel.Order)(nil),
		(*ordermodel.OrderItems)(nil),
		(*cartmodel.Cart)(nil),
		(*cartmodel.CartItems)(nil),
		(*deliverymodel.DeliveryPerson)(nil),
		(*deliverymodel.Deliveries)(nil),
	}); err != nil {
		log.Fatalf("order-svc: migrate: %s", err)
	}

	env := os.Getenv("APP_ENV")
	validate := validator.New()

	// gRPC client → User Service (:9091) for token validation on every auth'd request.
	userConn, err := grpc.NewClient("localhost:9091", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("order-svc: dial user-svc gRPC: %v", err)
	}
	defer userConn.Close()
	userGRPC := userpb.NewUserServiceClient(userConn)

	// gRPC client → Restaurant Service (:9092) for menu price lookups at checkout.
	restConn, err := grpc.NewClient("localhost:9092", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("order-svc: dial restaurant-svc gRPC: %v", err)
	}
	defer restConn.Close()
	restGRPC := restaurantpb.NewRestaurantServiceClient(restConn)

	// Auth middleware that calls User Service gRPC instead of reading JWT locally.
	authMiddleware := ordermiddleware.GRPCAuthMiddleware(userGRPC)

	// NATS for async order event publishing.
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://127.0.0.1:4222"
	}
	natServer, err := nats.NewNATS(natsURL)
	if err != nil {
		slog.Warn("order-svc: NATS unavailable — notifications disabled", "error", err)
	}

	wsClients := wsclients.NewRegistry()
	s := handler.NewServer(db, true)

	// Cart / Orders — uses gRPC price fetcher for checkout instead of DB join.
	cartService := cartsvc.NewCartServiceWithFetcher(db, env, natServer, &grpcPriceFetcher{client: restGRPC})
	carthandler.NewCartHandler(s, "/cart", cartService, []gin.HandlerFunc{authMiddleware}, validate)

	// Delivery
	deliveryService := deliverysvc.NewDeliveryService(db, env, natServer)
	deliveryhandler.NewDeliveryHandler(s, "/delivery", deliveryService, []gin.HandlerFunc{authMiddleware}, validate)

	// Reviews
	reviewService := reviewsvc.NewReviewService(db, env)
	reviewhandler.NewReviewProtectedHandler(s, "/review", reviewService, []gin.HandlerFunc{authMiddleware}, validate)

	// Announcements (public)
	announcementService := annsvc.NewAnnouncementService(db, env)
	annhandler.NewAnnouncementHandler(s, "/announcements", announcementService, []gin.HandlerFunc{authMiddleware}, validate)

	// Notifications / WebSocket
	notificationService := notifysvc.NewNotificationService(db, env, natServer)
	_ = notificationService.SubscribeNewOrders(wsClients)
	_ = notificationService.SubscribeOrderStatus(wsClients)
	notifyhandler.NewNotifyHandler(s, "/notify", notificationService, []gin.HandlerFunc{authMiddleware}, validate, wsClients)

	s.Gin.GET("/healthz", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"service": "order-svc"}) })

	slog.Info("order-svc: HTTP listening", "addr", ":8083")
	log.Fatal(s.Gin.Run(":8083"))
}
