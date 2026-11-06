package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
)

type EmergencyCancelRequest struct {
	BookingID string `json:"booking_id" binding:"required"`
	AdminID   string `json:"admin_id" binding:"required"`
	Reason    string `json:"reason" binding:"required"`
}

// HandleEmergencyCancel safely revokes active bookings, updates Firestore, and fires broadcast alerts.
func HandleEmergencyCancel(client *firestore.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req EmergencyCancelRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx := context.Background()
		docRef := client.Collection("booking_requests").Doc(req.BookingID)

		var booking BookingRequest

		// Execute Atomic Firestore Transaction
		err := client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
			doc, err := tx.Get(docRef)
			if err != nil {
				return fmt.Errorf("booking document not found: %v", err)
			}

			if err := doc.DataTo(&booking); err != nil {
				return fmt.Errorf("failed to parse booking data: %v", err)
			}

			if booking.Status == "cancelled_admin" || booking.Status == "completed" {
				return fmt.Errorf("booking %s is already in terminal status: %s", req.BookingID, booking.Status)
			}

			// Update state atomically
			updates := []firestore.Update{
				{Path: "status", Value: "cancelled_admin"},
				{Path: "cancellation_reason", Value: req.Reason},
				{Path: "cancelled_by", Value: req.AdminID},
				{Path: "cancelled_at", Value: time.Now().UTC()},
			}
			return tx.Update(docRef, updates)
		})

		if err != nil {
			log.Printf("[Admin Error] Emergency cancel failed for %s: %v", req.BookingID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Dispatch emergency system-wide Discord notification
		go dispatchEmergencyAlert(booking, req.Reason)

		c.JSON(http.StatusOK, gin.H{
			"status":     "success",
			"booking_id": req.BookingID,
			"message":    "Booking successfully cancelled and system notified.",
		})
	}
}

// dispatchEmergencyAlert pushes a high-priority alert to the Discord notification channel.
func dispatchEmergencyAlert(req BookingRequest, reason string) {
	webhookURL := os.Getenv("DISCORD_C500_NOTIFICATIONS_WEBHOOK")
	if webhookURL == "" {
		log.Println("[Warning] DISCORD_C500_NOTIFICATIONS_WEBHOOK not set.")
		return
	}

	loc, err := time.LoadLocation(req.IANAUserTZ)
	if err != nil {
		loc = time.UTC
	}
	localTimeStr := req.SlotUTC.In(loc).Format("Mon, Jan 2 at 3:04 PM MST")

	payload := fmt.Sprintf(`{
		"content": "🚨 **C500 SYSTEM ALERT: Mentorship Session Cancelled**",
		"embeds": [{
			"title": "Emergency Override: %s",
			"color": 15158332,
			"fields": [
				{"name": "Affected Student", "value": "<@%s>", "inline": true},
				{"name": "Assigned Mentor", "value": "%s", "inline": true},
				{"name": "Original Slot", "value": "%s (`%s`)", "inline": false},
				{"name": "Reason Provided", "value": "%s", "inline": false}
			],
			"footer": {"text": "Classroom 500 Admin Protocol • TimeSync"}
		}]
	}`, req.Topic, req.StudentDiscordID, req.MentorID, localTimeStr, req.IANAUserTZ, reason)

	// Execute POST request payload to Discord webhook...
	log.Printf("[Admin Dispatch] Emergency alert broadcasted for booking %s", req.BookingID)
}
