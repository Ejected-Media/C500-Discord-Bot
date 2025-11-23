This is where the magic happens. We turn that raw JSON data into a usable Go object and present the user with the Mentor Application Form.
To do this, we need to introduce Go's html/template engine. This allows us to inject the user's Discord data (like their username and avatar) directly into the HTML page before sending it to their browser.
1. The Data Structure
First, we need a Go struct that matches the JSON Discord sends us. Add this to your main.go file (outside of any function).
// User represents the data we get back from Discord
type DiscordUser struct {
    ID            string `json:"id"`
    Username      string `json:"username"`
    Discriminator string `json:"discriminator"` // The #1234 part (legacy but still used)
    GlobalName    string `json:"global_name"`   // The new display name
    Avatar        string `json:"avatar"`
    Email         string `json:"email"`
}

2. The HTML Template
Instead of hardcoding HTML strings inside the handler, we define a template. This template will have "holes" ({{.Username}}) that Go fills in.
Add this variable near the top of your file:
import "html/template"

// A simple HTML form styled for C500
const formHTML = `
<!DOCTYPE html>
<html>
<head>
    <title>C500 Mentor Application</title>
    <style>
        body { font-family: sans-serif; background: #2c2f33; color: white; display: flex; justify-content: center; padding: 50px; }
        .container { background: #23272a; padding: 40px; border-radius: 8px; width: 500px; box-shadow: 0 4px 15px rgba(0,0,0,0.5); }
        .profile { display: flex; align-items: center; margin-bottom: 20px; border-bottom: 1px solid #444; padding-bottom: 20px; }
        .avatar { width: 50px; height: 50px; border-radius: 50%; margin-right: 15px; }
        input, textarea { width: 100%; padding: 10px; margin-top: 5px; margin-bottom: 15px; background: #40444b; border: none; color: white; border-radius: 4px; }
        button { background: #7289da; color: white; padding: 12px 20px; border: none; border-radius: 4px; cursor: pointer; font-weight: bold; width: 100%; }
        button:hover { background: #5b6eae; }
        label { font-weight: bold; font-size: 0.9em; color: #b9bbbe; }
    </style>
</head>
<body>

<div class="container">
    <div class="profile">
        <img class="avatar" src="https://cdn.discordapp.com/avatars/{{.ID}}/{{.Avatar}}.png" alt="Avatar">
        <div>
            <h2>Welcome, {{.GlobalName}}</h2>
            <small>ID: {{.ID}}</small>
        </div>
    </div>

    <form action="/submit-application" method="POST">
        <input type="hidden" name="discord_id" value="{{.ID}}">
        <input type="hidden" name="discord_name" value="{{.GlobalName}}">

        <label>Twitch Channel URL</label>
        <input type="text" name="twitch_url" placeholder="https://twitch.tv/..." required>

        <label>Your CSS Philosophy (Why layout matters)</label>
        <textarea name="philosophy" rows="4" placeholder="I believe CSS Grid is..." required></textarea>

        <label>Link to your best CodePen / Work</label>
        <input type="text" name="work_link_1" placeholder="https://codepen.io/..." required>

        <button type="submit">Submit Application</button>
    </form>
</div>

</body>
</html>
`

3. The Updated Logic (main.go)
Now we update the handleCallback function to use the struct and the template.
Here is the complete, updated main.go combining everything.
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "html/template"
    "log"
    "net/http"
    "os"

    "golang.org/x/oauth2"
    "golang.org/x/oauth2/discord"
)

// --- CONFIG ---
// Replace these with your actual keys from the Developer Portal!
var conf = &oauth2.Config{
    ClientID:     "YOUR_CLIENT_ID",
    ClientSecret: "YOUR_CLIENT_SECRET", 
    RedirectURL:  "http://localhost:8080/auth/callback",
    Scopes:       []string{"identify", "email"},
    Endpoint:     discord.Endpoint,
}

// --- DATA MODELS ---
type DiscordUser struct {
    ID         string `json:"id"`
    Username   string `json:"username"`
    GlobalName string `json:"global_name"` 
    Avatar     string `json:"avatar"`
}

// --- TEMPLATE ---
// (Paste the 'const formHTML' string from section 2 here)
// For brevity, I am assuming the const is defined above.
var tmpl = template.Must(template.New("form").Parse(formHTML))

func main() {
    http.HandleFunc("/", handleHome)
    http.HandleFunc("/login", handleLogin)
    http.HandleFunc("/auth/callback", handleCallback)
    http.HandleFunc("/submit-application", handleSubmit) // New Handler!

    fmt.Println("C500 Auth Server running at http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleHome(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte(`<a href="/login">Login to Apply</a>`))
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
    url := conf.AuthCodeURL("random-state-string", oauth2.AccessTypeOnline)
    http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
    ctx := context.Background()

    // 1. Get Code
    code := r.URL.Query().Get("code")
    if code == "" {
        http.Error(w, "No code provided", http.StatusBadRequest)
        return
    }

    // 2. Exchange Code for Token
    token, err := conf.Exchange(ctx, code)
    if err != nil {
        http.Error(w, "Exchange failed", http.StatusInternalServerError)
        return
    }

    // 3. Fetch User Profile
    client := conf.Client(ctx, token)
    resp, err := client.Get("https://discord.com/api/users/@me")
    if err != nil {
        http.Error(w, "Fetch failed", http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    // 4. Parse JSON into Struct
    var user DiscordUser
    if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
        http.Error(w, "JSON parse failed", http.StatusInternalServerError)
        return
    }

    // 5. Render HTML Form with User Data
    // This injects user.ID and user.Avatar into the HTML
    if err := tmpl.Execute(w, user); err != nil {
        http.Error(w, "Template render failed", http.StatusInternalServerError)
    }
}

// Placeholder for Step 3
func handleSubmit(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    // Grab the form data
    philosophy := r.FormValue("philosophy")
    discordName := r.FormValue("discord_name")
    
    fmt.Fprintf(w, "Application Received for %s! Philosophy: %s", discordName, philosophy)
}

4. Test It
 * Run the code: go run main.go
 * Go to http://localhost:8080/login
 * Authorize with Discord.
Expected Result:
You should be redirected to a dark-mode page showing your own Discord Avatar and name at the top, followed by the form to input your Twitch URL and CSS philosophy.
Why this is cool for C500
You aren't just "asking" for their name. By rendering their avatar, you are proving to the user that C500 knows who they are. It builds trust and makes the onboarding feel professional and integrated.
Next Step:
Right now, the form just prints to the screen. To make this a real intake system, we need to Step 3: Save this form data to Firestore (and ideally ping your Discord server).
Ready for the Database connection?
