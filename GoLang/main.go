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
