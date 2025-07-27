// Ping command.
//
// Usage: /ping
package ping

import (
	"discordgo-bot/core/cmds"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

func handlecommand(
	cdata *cmds.CommandCreateData,
	parameters map[string]*discordgo.ApplicationCommandInteractionDataOption,
) {
	msg_desc := fmt.Sprintf("# Pong! 🏓\n-# %dms response time", cdata.Session.HeartbeatLatency().Milliseconds())

	embed := &discordgo.MessageEmbed{
		Title:       " ",
		Description: msg_desc,
		Color:       0x41aa0e,
	}

	err := fmt.Errorf("message & interaction creates are inactive")
	switch cdata.GetActive() {
	case cmds.CreateMessageType:
		_, err = cdata.Session.ChannelMessageSendComplex(
			cdata.Message.ChannelID,
			&discordgo.MessageSend{
				Embed: embed,
				Reference: &discordgo.MessageReference{
					MessageID: cdata.Message.ID,
					ChannelID: cdata.Message.ChannelID,
					GuildID:   cdata.Message.GuildID,
				},
				AllowedMentions: &discordgo.MessageAllowedMentions{
					RepliedUser: false,
				},
			},
		)
	case cmds.CreateInteractionType:
		err = cdata.Session.InteractionRespond(cdata.Interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed},
			},
		},
		)
	}
	if err != nil {
		log.Println(err)
	}
}

func init() {
	cmds.Register(cmds.CommandEntry{
		AppCommand: discordgo.ApplicationCommand{
			Name:        "ping",
			Description: "Pong! Responds with response latency.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "page",
					Description: "Page to display.",
					Type:        discordgo.ApplicationCommandOptionInteger,
					Required:    false,
				},
			},
		},
		HandleFunc: handlecommand,
	})
}
