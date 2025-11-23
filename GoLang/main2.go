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
