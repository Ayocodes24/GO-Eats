package main

import (
	"log"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func proxy(target string) gin.HandlerFunc {
	remote, err := url.Parse(target)
	if err != nil {
		log.Fatalf("gateway: invalid proxy target %q: %v", target, err)
	}
	rp := httputil.NewSingleHostReverseProxy(remote)
	return func(c *gin.Context) {
		rp.ServeHTTP(c.Writer, c.Request)
	}
}

func main() {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Route each path prefix to its owning microservice.
	r.Any("/user/*path",          proxy("http://localhost:8081"))
	r.Any("/restaurant/*path",    proxy("http://localhost:8082"))
	r.Any("/cart/*path",          proxy("http://localhost:8083"))
	r.Any("/delivery/*path",      proxy("http://localhost:8083"))
	r.Any("/review/*path",        proxy("http://localhost:8083"))
	r.Any("/announcements/*path", proxy("http://localhost:8083"))
	r.Any("/notify/*path",        proxy("http://localhost:8083"))

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "gateway"})
	})

	slog.Info("gateway: listening", "addr", ":8080")
	log.Fatal(r.Run(":8080"))
}
