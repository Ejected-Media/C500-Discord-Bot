// ... (imports)

// Global state to track the active session (Simple version)
// In production, use a database or map if managing multiple servers
var activeStageChannelID string

func handleOfficeHours(s *discordgo.Session, m *discordgo.MessageCreate) {
    // 1. Parse the Command
    args := strings.SplitN(m.Content, " ", 3)
    if len(args) < 2 {
        s.ChannelMessageSend(m.ChannelID, "Usage: `!officehours start [Topic]` or `!officehours end`")
        return
    }
    action := args[1]

    // 2. START SESSION
    if action == "start" {
        if activeStageChannelID != "" {
            s.ChannelMessageSend(m.ChannelID, "⚠️ A session is already active!")
            return
        }
        
        topic := "General Q&A"
        if len(args) > 2 {
            topic = args[2]
        }

        // A. Create the Stage Channel
        // We deny "Speak" for @everyone, ensuring only the Mentor (User) can talk initially.
        channel, err := s.GuildChannelCreateComplex(m.GuildID, discordgo.GuildChannelCreateData{
            Name: "🔴 LIVE: " + topic,
            Type: discordgo.ChannelTypeGuildStageVoice, 
            PermissionOverwrites: []*discordgo.PermissionOverwrite{
                {
                    ID:   m.GuildID, // @everyone role ID is same as Guild ID
                    Type: discordgo.PermissionOverwriteTypeRole,
                    Deny: discordgo.PermissionSpeak, // Students can only listen
                },
            },
        })

        if err != nil {
            s.ChannelMessageSend(m.ChannelID, "Error creating stage: "+err.Error())
            return
        }

        // B. "Open" the Stage (Create the Stage Instance)
        // This is what makes it appear at the top of the server list with the green "Live" badge.
        _, err = s.StageInstanceCreate(&discordgo.StageInstanceParams{
            ChannelID: channel.ID,
            Topic:     topic,
        })

        if err != nil {
            s.ChannelMessageSend(m.ChannelID, "Channel created, but failed to go LIVE: "+err.Error())
        }

        activeStageChannelID = channel.ID

        // C. Post the Announcement Embed
        embed := &discordgo.MessageEmbed{
            Title:       "🚨 Office Hours Starting Now!",
            Description: fmt.Sprintf("**Topic:** %s\n\nJoin the stage to listen in. Raise your hand ✋ to ask a question.", topic),
            Color:       0x5865F2, // Blurple
            Footer:      &discordgo.MessageEmbedFooter{Text: "Host: " + m.Author.Username},
        }
        s.ChannelMessageSendEmbed(m.ChannelID, embed)
    }

    // 3. END SESSION
    if action == "end" {
        if activeStageChannelID == "" {
            s.ChannelMessageSend(m.ChannelID, "No active session to end.")
            return
        }

        // Delete the channel (This automatically kills the Stage Instance)
        _, err := s.ChannelDelete(activeStageChannelID)
        if err != nil {
            s.ChannelMessageSend(m.ChannelID, "Failed to close stage: "+err.Error())
            return
        }

        s.ChannelMessageSend(m.ChannelID, "✅ Office Hours concluded. Thanks for joining!")
        activeStageChannelID = ""
    }
}
