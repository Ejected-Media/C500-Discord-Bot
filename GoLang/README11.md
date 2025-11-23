This is the foundational security layer for C500.tech. Implementing "Login with Discord" allows you to verify that every mentor applicant is a real person with a real Discord account, without you needing to manage passwords or store sensitive email data.
Here is the implementation plan for the Go OAuth2 Handshake.
1. The Discord Developer Setup (Prerequisite)
Before writing code, you need to register the app to get your "Keys" to the castle.
 * Go to the Discord Developer Portal.
 * Create a New Application called "Classroom 500".
 * Go to the OAuth2 tab.
 * Add a Redirect Method: http://localhost:8080/auth/callback
   * Note: When you deploy to C500.tech later, you will add https://c500.tech/auth/callback here too.
 * Copy your Client ID and Client Secret.
2. The Project Setup
Create a new folder for this project (e.g., c500-onboarding) and initialize it:
mkdir c500-onboarding
cd c500-onboarding
go mod init c500-onboarding
go get golang.org/x/oauth2

3. The Go Server (main.go)
This single file handles the entire secure handshake.
How it works:
 * /login: The user clicks this. We construct a secure URL and kick them over to Discord.
 * Discord: Asks the user "Do you allow C500 to access your username?"
 * /auth/callback: Discord kicks them back to us with a secret code.
 * Exchange: We swap that code for an access_token (the key) and use it to get their User ID.
<!-- end list -->
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"

    "golang.org/x/oauth2"
    "golang.org/x/oauth2/discord"
)

// CONFIGURATION
// In production, fetch these from os.Getenv() for security
var conf = &oauth2.Config{
    ClientID:     "YOUR_CLIENT_ID_HERE",     // <-- PASTE FROM PORTAL
    ClientSecret: "YOUR_CLIENT_SECRET_HERE", // <-- PASTE FROM PORTAL
    RedirectURL:  "http://localhost:8080/auth/callback",
    Scopes:       []string{"identify", "email"},
    Endpoint:     discord.Endpoint,
}

func main() {
    http.HandleFunc("/", handleHome)
    http.HandleFunc("/login", handleLogin)
    http.HandleFunc("/auth/callback", handleCallback)

    fmt.Println("C500 Auth Server running at http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

// 1. THE LANDING PAGE
func handleHome(w http.ResponseWriter, r *http.Request) {
    html := `<html><body>
    <h1>Welcome to Classroom 500</h1>
    <p>Mentor Onboarding Portal</p>
    <a href="/login"><button>Login with Discord</button></a>
    </body></html>`
    w.Write([]byte(html))
}

// 2. THE REDIRECT (Send them to Discord)
func handleLogin(w http.ResponseWriter, r *http.Request) {
    // We generate the URL that asks Discord for permission
    url := conf.AuthCodeURL("state-token", oauth2.AccessTypeOnline)
    http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// 3. THE HANDSHAKE (They came back!)
func handleCallback(w http.ResponseWriter, r *http.Request) {
    ctx := context.Background()

    // A. Get the "code" Discord sent us in the URL query params
    code := r.URL.Query().Get("code")
    if code == "" {
        http.Error(w, "Code not found", http.StatusBadRequest)
        return
    }

    // B. Exchange the "code" for an actual "Access Token"
    // This happens securely between our server and Discord's server
    token, err := conf.Exchange(ctx, code)
    if err != nil {
        http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // C. Use the Token to ask Discord "Who is this user?"
    client := conf.Client(ctx, token)
    resp, err := client.Get("https://discord.com/api/users/@me")
    if err != nil {
        http.Error(w, "Failed to get user info", http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    // D. Read and Display the User Data (For now)
    userData, _ := io.ReadAll(resp.Body)
    
    // In the next step, we will Parse this JSON and save it to a Cookie/Session
    w.Header().Set("Content-Type", "application/json")
    w.Write(userData)
}

4. Running the Test
 * Paste your actual Client ID and Secret into the var conf block.
 * Run go run main.go.
 * Open http://localhost:8080 in your browser.
 * Click Login with Discord.
If it works, you should be redirected back to your localhost, and you will see a raw JSON dump of your own Discord profile (ID, Username, Avatar hash) on the screen.
Critical Concept: The "State" Token
In the code above, I used "state-token" as a placeholder. In a production security environment, you generate a random string, save it in a cookie, and check it when they come back. This prevents CSRF (Cross-Site Request Forgery) attacks.
 * Do you want to add that security layer now, or keep it simple for the prototype?
Next Step
Once you confirm you can see your own JSON data on the screen, we need to parse that JSON into a struct and render the Mentor Application Form.
Would you like to proceed to Step 2: Parsing User Data & Rendering the Form?
