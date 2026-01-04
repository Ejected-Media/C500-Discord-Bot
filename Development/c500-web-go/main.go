package main

import (
	"log"
	"github.com/gin-gonic/gin"
	"c500-web-go/internal/config"
	"c500-web-go/internal/clients"
	"c500-web-go/internal/handlers"
)

func main() {
	// Load config and init Core API client
	cfg := config.Load()
	coreClient := clients.NewCoreAPIClient(cfg.CoreAPIURL, cfg.InternalAPIKey)

	// Init Handlers
	profileHandler := handlers.NewProfileHandler(coreClient)
	// staticHandler := ...

	r := gin.Default()

	// 1. Load HTML Templates (glob pattern finds them all)
	r.LoadHTMLGlob("templates/**/*")

	// 2. Serve Static Assets directly
	// Requests to /assets/css/main.css will serve the file from local ./assets/css/main.css
	r.Static("/assets", "./assets")

	// 3. Define Routes
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "home.html", gin.H{"Title": "Home"})
	})

	// The public profile route
	r.GET("/builder/:username", profileHandler.GetBuilderProfile)

	// Stripe redirect routes
	r.GET("/success", func(c *gin.Context) { /* render success.html */ })
	r.GET("/cancel", func(c *gin.Context) { /* render cancel.html */ })

	log.Println("C500 Web Server starting...")
	r.Run(":" + cfg.Port)
}
