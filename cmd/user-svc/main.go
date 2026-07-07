package main

import (
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/Ayocodes24/GO-Eats/pkg/database"
	"github.com/Ayocodes24/GO-Eats/pkg/database/models/user"
	"github.com/Ayocodes24/GO-Eats/pkg/grpc/userserver"
	"github.com/Ayocodes24/GO-Eats/pkg/handler"
	userhandler "github.com/Ayocodes24/GO-Eats/pkg/handler/user"
	userpb "github.com/Ayocodes24/GO-Eats/pkg/proto/user"
	usersvc "github.com/Ayocodes24/GO-Eats/pkg/service/user"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	if err := godotenv.Load(".env.user"); err != nil {
		log.Fatal("Error loading .env.user")
	}

	db := database.New()
	if err := db.MigrateModels([]interface{}{(*user.User)(nil)}); err != nil {
		log.Fatalf("user-svc: migrate: %s", err)
	}

	// gRPC server on :9091 — Order Service calls ValidateToken here.
	go func() {
		lis, err := net.Listen("tcp", ":9091")
		if err != nil {
			log.Fatalf("user-svc: gRPC listen: %v", err)
		}
		grpcSrv := grpc.NewServer()
		userpb.RegisterUserServiceServer(grpcSrv, userserver.New())
		slog.Info("user-svc: gRPC listening", "addr", ":9091")
		if err := grpcSrv.Serve(lis); err != nil {
			log.Fatalf("user-svc: gRPC serve: %v", err)
		}
	}()

	// HTTP server on :8081
	s := handler.NewServer(db, true)
	env := os.Getenv("APP_ENV")
	userhandler.NewUserHandler(s, "/user", usersvc.NewUserService(db, env), validator.New())
	s.Gin.GET("/healthz", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"service": "user-svc"}) })

	slog.Info("user-svc: HTTP listening", "addr", ":8081")
	log.Fatal(s.Gin.Run(":8081"))
}
