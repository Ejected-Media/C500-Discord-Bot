This is the final piece of the puzzle. We are going to take the data from that HTML form and write it directly to your Google Cloud database.
As a bonus, I will also add a Discord Webhook integration. This means the moment someone applies, your private admin channel on Discord will get a ping saying: "New Mentor Application from @User!"
1. The Data Model (Application)
We need a struct to define what a "Mentor Application" looks like in the database. Add this to your main.go.
import "time"

type MentorApplication struct {
    DiscordID    string    `firestore:"discord_id"`
    DiscordName  string    `firestore:"discord_name"`
    TwitchURL    string    `firestore:"twitch_url"`
    Philosophy   string    `firestore:"philosophy"`
    WorkLink     string    `firestore:"work_link"`
    Status       string    `firestore:"status"` // e.g., "pending", "approved"
    SubmittedAt  time.Time `firestore:"submitted_at"`
}

2. The Database Connection
We need to initialize Firestore when the app starts. This is almost identical to the Drift Miner backend.
Add these imports:
import (
    "cloud.google.com/go/firestore"
    firebase "firebase.google.com/go/v4"
    "google.golang.org/api/option"
)

Add this global variable:
var fsClient *firestore.Client

Add this initFirestore function:
func initFirestore() {
    ctx := context.Background()
    conf := &firebase.Config{ProjectID: "YOUR_PROJECT_ID"} // e.g. c500-tech

    // Note: On App Engine, you don't need option.WithCredentialsFile
    // For local dev, point to your serviceAccountKey.json
    opt := option.WithCredentialsFile("serviceAccountKey.json")
    
    app, err := firebase.NewApp(ctx, conf, opt)
    if err != nil {
        log.Fatalf("error initializing app: %v", err)
    }

    client, err := app.Firestore(ctx)
    if err != nil {
        log.Fatalf("error getting firestore client: %v", err)
    }
    fsClient = client
}

Call initFirestore() inside your main() function before starting the server.
3. The "Discord Ping" (Webhook)
This is the "secret sauce." Instead of checking the database manually, let the database notify you.
 * Go to your Discord Server -> Server Settings -> Integrations -> Webhooks.
 * Create a new Webhook (call it "C500 Bot") and copy the URL.
Add this helper function:
import "bytes"

func sendDiscordNotification(app MentorApplication) {
    webhookURL := "YOUR_WEBHOOK_URL_HERE"

    // Format the message using Discord's JSON structure
    // We can mention roles using <@&ROLE_ID> if you want to ping Admins
    jsonBody := []byte(fmt.Sprintf(`{
        "content": "🚨 **New Mentor Application!**",
        "embeds": [{
            "title": "%s applied to be a mentor",
            "color": 5763719,
            "fields": [
                {"name": "Philosophy", "value": "%s"},
                {"name": "Twitch", "value": "%s"},
                {"name": "Link", "value": "%s"}
            ]
        }]
    }`, app.DiscordName, app.Philosophy, app.TwitchURL, app.WorkLink))

    http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonBody))
}

4. The Final handleSubmit
Update your handleSubmit function to bring it all together.
func handleSubmit(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    ctx := context.Background()

    // 1. Construct the Application Object
    app := MentorApplication{
        DiscordID:   r.FormValue("discord_id"),
        DiscordName: r.FormValue("discord_name"),
        TwitchURL:   r.FormValue("twitch_url"),
        Philosophy:  r.FormValue("philosophy"),
        WorkLink:    r.FormValue("work_link_1"),
        Status:      "pending",
        SubmittedAt: time.Now(),
    }

    // 2. Save to Firestore
    // We use .Add() to let Firestore generate a unique ID
    _, _, err := fsClient.Collection("mentors").Add(ctx, app)
    if err != nil {
        log.Printf("Failed to save application: %v", err)
        http.Error(w, "Database Error", http.StatusInternalServerError)
        return
    }

    // 3. Ping Discord (Fire and Forget)
    go sendDiscordNotification(app)

    // 4. Show Success Page
    w.Write([]byte(`
        <html>
        <body style="background:#2c2f33; color:white; text-align:center; padding-top:50px; font-family:sans-serif;">
            <h1>Application Received! 🚀</h1>
            <p>Thanks, ` + app.DiscordName + `. We have pinged the admins.</p>
            <p>You can close this window.</p>
        </body>
        </html>
    `))
}

How to Deploy this to C500.tech
Since you are using Google App Engine, deploying this is incredibly simple.
 * Create an app.yaml file in the root:
   runtime: go122
env_variables:
  DISCORD_CLIENT_ID: "your_id"
  DISCORD_CLIENT_SECRET: "your_secret"
  DISCORD_WEBHOOK_URL: "your_webhook"
  GOOGLE_CLOUD_PROJECT: "c500-tech"

 * Run command: gcloud app deploy
Summary
You now have a fully automated "Intake Machine":
 * User logs in securely via Discord.
 * User fills out the Go Template Form.
 * Data is saved to Firestore.
 * You get a notification in Discord Chat.
This effectively validates your C500.tech domain as a functional part of your ecosystem!
Would you like to create a repository for this one as well, or do you want to discuss the "CSS Layout" curriculum content that these mentors will be teaching?
