package main

import (
	"github.com/Ayocodes24/GO-Eats/cmd/api/middleware"
	"github.com/Ayocodes24/GO-Eats/pkg/abstract/config"
	"github.com/Ayocodes24/GO-Eats/pkg/database"
	"github.com/Ayocodes24/GO-Eats/pkg/handler/user"
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
}
