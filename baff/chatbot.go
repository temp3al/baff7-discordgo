package baff

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var (
	CHAT_CHANNEL   string = "1388969389872906290"
	chatlog_latest string
)

func handle_chatbot_read_message(s *discordgo.Session, m *discordgo.MessageCreate) {
	// follow CHAT_CHANNEL ID only & ignore our own messages
	if m.Author.ID == s.State.User.ID || m.ChannelID != CHAT_CHANNEL {
		return
	}
	chatlog_latest += fmt.Sprintf(
		"**%v:** %s\n",
		m.Author.DisplayName(),
		m.Content,
	)
	s.ChannelMessageSend(m.ChannelID, chatlog_latest)
}

func announce_publish() {

}

func ollama_request() {

}

// Send a discord message in the provided channel containing our response.
func send_message() {
}
