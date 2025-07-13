// Ping command.
//
// Usage: /ping
package ping

import (
	"discordgo-bot/core"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func help_command(session *discordgo.Session, message *discordgo.MessageCreate) error {
	desc_msg := fmt.Sprintf("# Pong! 🏓\n-# %dms response time", session.HeartbeatLatency().Milliseconds())

	_, err := session.ChannelMessageSendComplex(
		message.ChannelID,
		&discordgo.MessageSend{
			Embed: &discordgo.MessageEmbed{
				Title:       " ",
				Description: desc_msg,
				Color:       0x41aa0e,
			},
			Reference: &discordgo.MessageReference{
				MessageID: message.ID,
				ChannelID: message.ChannelID,
				GuildID:   message.GuildID,
			},
			AllowedMentions: &discordgo.MessageAllowedMentions{
				RepliedUser: true,
			},
		},
	)
	return err
}

func init() {
	core.RegisterCommand(
		core.Command{ // create command
			Name:        "ping",
			Description: "Pong! (Gets response latency.)",
			// chat message handle
			HandlerChat: func(session *discordgo.Session, message *discordgo.MessageCreate) error {
				return help_command(session, message)
			},
			// slash message handle
			HandlerSlash: func(session *discordgo.Session, interaction *discordgo.InteractionCreate) error {
				return nil
			},
		},
	)
}
