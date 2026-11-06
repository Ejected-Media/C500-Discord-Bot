package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	
	// Import your internal handlers and database packages
	// "c500-timesync-backend/internal/handlers"
)

func main() {
	ctx := context.Background()

	// 1. Retrieve GCP Project ID from environment (falls back to ADC project discovery)
	projectID := os.Getenv("FIRESTORE_PROJECT_ID")
	if projectID == "" {
		log.Println("[Warning] FIRESTORE_PROJECT_ID not set. Relying on default credential chain...")
	}

	// 2. Initialize Firebase Firestore Client
	firestoreClient, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("[Fatal] Failed to initialize Firestore client: %v", err)
	}
	defer firestoreClient.Close()
	log.Println("[Database] Connected to Firestore successfully.")

	// 3. Set up Gin Router
	router := gin.Default()

	// Global Middleware (CORS, Logging, Recovery)
	router.Use(gin.Recovery())

	// Health Check Endpoint (for Cloud Run / App Engine health probes)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "c500-timesync-backend",
		})
	})

	// API Route Group
	apiGroup := router.Group("/api")
	{
		// Mentorship booking creation endpoint (invoked by Discord bot / web portal)
		// apiGroup.POST("/bookings", handlers.HandleCreateBooking(firestoreClient))
	}

	// Cloud Scheduler / Task Queue Endpoint Group
	taskGroup := router.Group("/tasks")
	{
		// Periodic cron triggered every 5 minutes by Google Cloud Scheduler
		// taskGroup.POST("/check-reminders", handlers.HandleCheckReminders(firestoreClient))
	}

	// Admin Override Endpoint Group
	adminGroup := router.Group("/admin")
	{
		// Secure admin emergency cancellation protocol
		// adminGroup.POST("/emergency-cancel", handlers.HandleEmergencyCancel(firestoreClient))
	}

	// 4. Resolve Port and Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("[Server] Starting C500 TimeSync backend engine on port %s...", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("[Fatal] Server terminated unexpectedly: %v", err)
	}
}
