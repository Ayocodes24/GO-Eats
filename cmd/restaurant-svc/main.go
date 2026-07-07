package main

import (
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/Ayocodes24/GO-Eats/pkg/database"
	restaurantmodel "github.com/Ayocodes24/GO-Eats/pkg/database/models/restaurant"
	"github.com/Ayocodes24/GO-Eats/pkg/grpc/restaurantserver"
	"github.com/Ayocodes24/GO-Eats/pkg/handler"
	restauranthandler "github.com/Ayocodes24/GO-Eats/pkg/handler/restaurant"
	restaurantpb "github.com/Ayocodes24/GO-Eats/pkg/proto/restaurant"
	restaurantsvc "github.com/Ayocodes24/GO-Eats/pkg/service/restaurant"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	if err := godotenv.Load(".env.restaurant"); err != nil {
		log.Fatal("Error loading .env.restaurant")
	}

	db := database.New()
	if err := db.MigrateModels([]interface{}{
		(*restaurantmodel.Restaurant)(nil),
		(*restaurantmodel.MenuItem)(nil),
	}); err != nil {
		log.Fatalf("restaurant-svc: migrate: %s", err)
	}

	env := os.Getenv("APP_ENV")
	svc := restaurantsvc.NewRestaurantService(db, env)

	// gRPC server on :9092 — Order Service calls GetMenuItem here.
	go func() {
		lis, err := net.Listen("tcp", ":9092")
		if err != nil {
			log.Fatalf("restaurant-svc: gRPC listen: %v", err)
		}
		grpcSrv := grpc.NewServer()
		restaurantpb.RegisterRestaurantServiceServer(grpcSrv, restaurantserver.New(svc))
		slog.Info("restaurant-svc: gRPC listening", "addr", ":9092")
		if err := grpcSrv.Serve(lis); err != nil {
			log.Fatalf("restaurant-svc: gRPC serve: %v", err)
		}
	}()

	// HTTP server on :8082
	s := handler.NewServer(db, true)
	restauranthandler.NewRestaurantHandler(s, "/restaurant", svc)
	s.Gin.GET("/healthz", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"service": "restaurant-svc"}) })

	slog.Info("restaurant-svc: HTTP listening", "addr", ":8082")
	log.Fatal(s.Gin.Run(":8082"))
}
