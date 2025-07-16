package main

import (
	"github.com/Ayocodes24/GO-Eats/pkg/abstract/config"
	"github.com/Ayocodes24/GO-Eats/pkg/database"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load .env file
	config.LoadEnv()

	// Connect to DB
	db := database.New()
	defer db.Close()

	// Run auto-migration for tables
	if err := db.Migrate(); err != nil {
		panic("DB Migration failed: " + err.Error())
	}

	// Setup basic route
	r := gin.Default()
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Start server
	r.Run(":8080")
}
