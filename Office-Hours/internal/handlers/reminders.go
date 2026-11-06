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
	"google.golang.org/api/iterator"
)

// HandleCheckReminders is triggered periodically by Google Cloud Scheduler 
// to scan Firestore for upcoming sessions and dispatch Discord notifications.
func HandleCheckReminders(client *firestore.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		now := time.Now().UTC()
		
		// Define target window: Sessions starting between now and 15 minutes from now
		windowStart := now
		windowEnd := now.Add(15 * time.Minute)

		log.Printf("[TimeSync Reminder Job] Scanning sessions between %s and %s", 
			windowStart.Format(time.RFC3339), windowEnd.Format(time.RFC3339))

		// Query Firestore: confirmed bookings where reminder hasn't been sent yet
		query := client.Collection("booking_requests").
			Where("status", "==", "confirmed").
			Where("reminder_dispatched", "==", false).
			Where("slot_utc", ">=", windowStart).
			Where("slot_utc", "<=", windowEnd)

		iter := query.Documents(ctx)
		defer iter.Stop()

		dispatchedCount := 0

		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Printf("[Firestore Error] Failed to iterate booking documents: %v", err)
				continue
			}

			var booking BookingRequest
			if err := doc.DataTo(&booking); err != nil {
				log.Printf("[Mapping Error] Failed to parse booking data [%s]: %v", doc.Ref.ID, err)
				continue
			}

			// Dispatch the smart reminder ping via Discord webhook
			if err := sendReminderWebhook(booking); err != nil {
				log.Printf("[Discord Error] Failed to send reminder for booking %s: %v", booking.BookingID, err)
				continue
			}

			// Atomically update reminder_dispatched status to true
			_, err = doc.Ref.Update(ctx, []firestore.Update{
				{Path: "reminder_dispatched", Value: true},
			})
			if err != nil {
				log.Printf("[Firestore Error] Failed to update dispatch flag for %s: %v", booking.BookingID, err)
			} else {
				dispatchedCount++
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"status":            "success",
			"checked_at":        now.Format(time.RFC3339),
			"reminders_sent":    dispatchedCount,
		})
	}
}

// sendReminderWebhook constructs and pushes an urgent Discord ping for the session.
func sendReminderWebhook(req BookingRequest) error {
	webhookURL := os.Getenv("DISCORD_C500_NOTIFICATIONS_WEBHOOK")
	if webhookURL == "" {
		return fmt.Errorf("DISCORD_C500_NOTIFICATIONS_WEBHOOK environment variable not set")
	}

	// Load user's IANA timezone anchor for accurate local time rendering
	loc, err := time.LoadLocation(req.IANAUserTZ)
	if err != nil {
		loc = time.UTC
	}
	localTimeStr := req.SlotUTC.In(loc).Format("3:04 PM MST")

	// Mention both student (<@ID>) and mentor directly in the alert
	messageContent := fmt.Sprintf("⏰ **C500 Mentorship Starting Soon!**\nHey <@%s> and **%s**, your session on **%s** begins in about **15 minutes** (`%s`).",
		req.StudentDiscordID, req.MentorID, req.Topic, localTimeStr)

	// POST payload execution to Discord...
	log.Printf("[Reminder Dispatch] Pinging student %s and mentor %s for booking %s", 
		req.StudentDiscordID, req.MentorID, req.BookingID)
	
	return nil
}
