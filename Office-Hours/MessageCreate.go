func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
    if m.Author.ID == s.State.User.ID { return }

    // ... other commands ...

    if strings.HasPrefix(m.Content, "!officehours") {
        handleOfficeHours(s, m)
    }
}
