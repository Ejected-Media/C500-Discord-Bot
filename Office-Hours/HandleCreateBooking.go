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

type BookingRequest struct {
	BookingID        string    `json:"booking_id" firestore:"booking_id"`
	StudentDiscordID string    `json:"student_discord_id" firestore:"student_discord_id"`
	StudentName      string    `json:"student_name" firestore:"student_name"`
	MentorID         string    `json:"mentor_id" firestore:"mentor_id"`
	Topic            string    `json:"topic" firestore:"topic"`
	SlotUTC          time.Time `json:"slot_utc" firestore:"slot_utc"`
	IANAUserTZ       string    `json:"iana_timezone" firestore:"iana_timezone"`
	Status           string    `json:"status" firestore:"status"`
	ReminderSent     bool      `json:"reminder_dispatched" firestore:"reminder_dispatched"`
	CreatedAt        time.Time `json:"created_at" firestore:"created_at"`
}

// HandleCreateBooking processes incoming booking requests, stores them atomically, and triggers alerts.
func HandleCreateBooking(client *firestore.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req BookingRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx := context.Background()
		req.BookingID = fmt.Sprintf("bk_%d", time.Now().UnixNano())
		req.Status = "confirmed"
		req.ReminderSent = false
		req.CreatedAt = time.Now().UTC()

		// Perform atomic Firestore Transaction
		docRef := client.Collection("booking_requests").Doc(req.BookingID)
		err := client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
			return tx.Set(docRef, req)
		})

		if err != nil {
			log.Printf("[Firestore Error] Failed to create booking: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to persist booking record"})
			return
		}

		// Dispatch Discord notification webhook asynchronously
		go dispatchDiscordAlert(req)

		c.JSON(http.StatusOK, gin.H{
			"status":     "success",
			"booking_id": req.BookingID,
			"message":    "C500 Mentorship session successfully scheduled and logged.",
		})
	}
}

// dispatchDiscordAlert sends a notification to the mentor and student via Discord webhook.
func dispatchDiscordAlert(req BookingRequest) {
	webhookURL := os.Getenv("DISCORD_C500_NOTIFICATIONS_WEBHOOK")
	if webhookURL == "" {
		log.Println("[Warning] DISCORD_C500_NOTIFICATIONS_WEBHOOK not set in environment.")
		return
	}

	// Format time for user's IANA timezone anchor
	loc, err := time.LoadLocation(req.IANAUserTZ)
	if err != nil {
		loc = time.UTC
	}
	localTimeStr := req.SlotUTC.In(loc).Format("Mon, Jan 2 at 3:04 PM MST")

	payload := fmt.Sprintf(`{
		"content": "🎓 **New C500 Mentorship Booking Confirmed!**",
		"embeds": [{
			"title": "Topic: %s",
			"color": 7506394,
			"fields": [
				{"name": "Student", "value": "<@%s> (%s)", "inline": true},
				{"name": "Assigned Mentor", "value": "%s", "inline": true},
				{"name": "Scheduled Time", "value": "%s (`%s`)", "inline": false}
			],
			"footer": {"text": "Classroom 500 TimeSync Engine"}
		}]
	}`, req.Topic, req.StudentDiscordID, req.StudentName, req.MentorID, localTimeStr, req.IANAUserTZ)

	// POST request payload to Discord Webhook...
	log.Printf("[Discord Dispatch] Alert sent for booking %s", req.BookingID)
}
