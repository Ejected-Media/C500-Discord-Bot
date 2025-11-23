That is actually a brilliant insight. Because Firefox Focus is stateless (it wipes everything when you close it), it effectively forces a "Hard Refresh" every single time you use it.
For a web developer, Caching is often the enemy. You change a CSS color, reload the page, and nothing happens because the browser is "remembering" the old file. Using Firefox Focus as your "Previewer" guarantees you are seeing the site exactly as a brand-new visitor would, with zero cached data interfering.
I will add that to the "Mobile Studio" definition as the "Clean Slate Viewer."
🤖 Phase 3: The C500 Discord Bot (in Go)
Now, let's build the machine that runs the school.
Since you are using Go, we will use the industry-standard library: discordgo.
This bot needs to do two specific things to support your new "Tuition/Mentor" model:
 * The Intake Handler: When the Mentor Portal sends that webhook, the Bot should add "Approve / Deny" buttons to the message so you can handle it with one click.
 * The Office Hours Manager: A command to announce and manage live sessions.
1. Project Setup
Create a new directory (or add to your repo) called c500-bot.
mkdir c500-bot
cd c500-bot
go mod init c500-bot
go get github.com/bwmarrin/discordgo
go get github.com/joho/godotenv // For managing .env keys

2. The Bot Skeleton (main.go)
Here is the clean, modular structure. It connects to Discord and waits for the "Ready" signal.
package main

import (
    "flag"
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"

    "github.com/bwmarrin/discordgo"
    "github.com/joho/godotenv"
)

// Variables used for command line parameters
var (
    Token string
)

func init() {
    // Load .env file if it exists (good for local dev)
    godotenv.Load()
    
    flag.StringVar(&Token, "t", os.Getenv("DISCORD_BOT_TOKEN"), "Bot Token")
    flag.Parse()
}

func main() {
    if Token == "" {
        log.Fatal("No token provided. Please set DISCORD_BOT_TOKEN.")
    }

    // 1. Create a new Discord session
    dg, err := discordgo.New("Bot " + Token)
    if err != nil {
        log.Fatalf("error creating Discord session: %v", err)
    }

    // 2. Register Handlers
    // We will add these functions in the next step
    dg.AddHandler(ready)
    dg.AddHandler(messageCreate)
    dg.AddHandler(interactionCreate) // For button clicks

    // 3. Open Connection
    // We only need the 'Intents' to read messages and see guild members
    dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuilds

    err = dg.Open()
    if err != nil {
        log.Fatalf("error opening connection: %v", err)
    }

    fmt.Println("C500 Bot is now running. Press CTRL-C to exit.")

    // 4. Wait for Termination Signal (Keep app running)
    sc := make(chan os.Signal, 1)
    signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
    <-sc

    // 5. Cleanly Close
    dg.Close()
}

// --- HANDLERS ---

func ready(s *discordgo.Session, event *discordgo.Ready) {
    // Set the "Playing" status
    s.UpdateGameStatus(0, "CSS Layout | Level 00")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
    // Ignore all messages created by the bot itself
    if m.Author.ID == s.State.User.ID {
        return
    }

    // Simple Ping Test
    if m.Content == "!ping" {
        s.ChannelMessageSend(m.ChannelID, "Pong! C500 Systems Online.")
    }
}

3. The "Interactive" Feature (Approve/Deny)
This is the cool part. When your Web Portal sends a webhook, it's just text. But the Bot can post a message with Buttons.
Let's create a command !review_latest (or trigger this automatically via Firestore listeners later) that posts an interactive card for a mentor application.
Add this function to main.go (inside messageCreate or as a separate handler):
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
    if m.Author.ID == s.State.User.ID { return }

    // Simulate a Moderator reviewing an application
    if m.Content == "!test_intake" {
        
        // Define the "Approve" and "Deny" buttons
        approveBtn := discordgo.Button{
            Label:    "Approve Mentor",
            Style:    discordgo.SuccessButton, // Green
            CustomID: "approve_mentor_123",    // We will parse this ID later
        }

        denyBtn := discordgo.Button{
            Label:    "Deny",
            Style:    discordgo.DangerButton,  // Red
            CustomID: "deny_mentor_123",
        }

        // Create the Message with the Buttons attached
        msg := &discordgo.MessageSend{
            Content: "📄 **Application Review**\n**User:** @PixelNinja\n**Level:** 03 (Grid Master)",
            Components: []discordgo.MessageComponent{
                discordgo.ActionsRow{
                    Components: []discordgo.MessageComponent{approveBtn, denyBtn},
                },
            },
        }

        s.ChannelMessageSendComplex(m.ChannelID, msg)
    }
}

// Handle the Button Click
func interactionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
    if i.Type == discordgo.InteractionMessageComponent {
        
        customID := i.MessageComponentData().CustomID
        
        if customID == "approve_mentor_123" {
            // 1. Update the message to remove buttons (so it can't be clicked twice)
            s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
                Type: discordgo.InteractionResponseUpdateMessage,
                Data: &discordgo.InteractionResponseData{
                    Content: "✅ **Mentor Approved!** Role assigned.",
                    Components: []discordgo.MessageComponent{}, // Empty list removes buttons
                },
            })
            
            // 2. TODO: Assign the Discord Role here
            // s.GuildMemberRoleAdd(i.GuildID, "MEMBER_ID", "MENTOR_ROLE_ID")
        }
    }
}

4. Deployment Strategy
Since this is a "Long-Running Process" (it needs to stay connected to the socket), it is slightly different from the Web Portal we built earlier.
 * Web Portal: App Engine Standard (Scales to zero, HTTP only).
 * Discord Bot: App Engine Flexible OR (Better/Cheaper) Google Compute Engine (e2-micro).
Since you want to keep costs low, running this on a tiny Compute Engine instance (which is free tier eligible in many regions) is the "Independent" way to do it. You just SSH in, run the binary, and walk away.
Does this Bot structure look good to you? If so, we can write the specific logic for the Office Hours command next.
