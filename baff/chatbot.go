// Package providing chatbot communication via ollama.
package baff

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var (
	CHAT_CHANNEL string = "1388969389872906290"
	OLLAMA_MODEL string

	chatlog_latest string
	chat_data      map[string]channel_chatdata // "ChannelID": channel_chatdata
)

type channel_chatdata struct {
}

// Perform a request to our running ollama server.
func ollama_api_perform_request() {

}

func handle_readmessage_chatbot(s *discordgo.Session, m *discordgo.MessageCreate) {
	// follow CHAT_CHANNEL ID only & ignore our own messages
	if m.Author.ID == s.State.User.ID || m.ChannelID != CHAT_CHANNEL {
		return
	}
	chatlog_latest += fmt.Sprintf(
		"**%v:** %s\n",
		m.Author.DisplayName(),
		m.Content,
	)
}

//
//
//

// Add our newly received message into our query string, start or reset a timer before we start our request.
func chatbot_process_message() {

}

// Use our final message to send a request to our ollama server.
func chatbot_start_request() {

}

// Send a chat message to the desired channel containing our ollama response.
func chatbot_finalize() {

}
