// Help command.
//
// Usage: /help [page]
package help

import (
	"discordgo-bot/core/cmds"
	"fmt"
	"log"
	"math"

	"github.com/bwmarrin/discordgo"
)

var (
	// Amount of commands to be rendered per command call.
	//
	// Higher is more information at the cost of chat visibility.
	commands_per_page int = 11
)

func handlecommand(
	cdata *cmds.CommandCreateData,
	parameters map[string]*discordgo.ApplicationCommandInteractionDataOption,
) {
	cmd_prefix := "/" // core.Prefixes[0]
	cmd_map := cmds.GetCommandEntries()

	page_max := int(math.Ceil(
		float64(len(cmd_map)) /
			float64(commands_per_page),
	))
	page_active := 1
	// pull our page parameter
	if pmt, ok := parameters["page"]; ok {
		page_active = max(1, min(page_max, int(pmt.IntValue())))
	}

	str_title := "Command List"
	// generate description from core
	str_desc := ""
	i := 0 // ranging maps doesnt return a len variable so...
	for _, cmd := range cmd_map {
		i++
		if i < commands_per_page*(page_active-1) { // offset our command list using help_page & commands_per_page
			continue
		} else if i+1 > commands_per_page*page_active { // stop generating if we go past our cmds limit.
			break
		}
		// generate usage line using command options
		str_pmtrs := ""
		cmd_options := cmd.AppCommand.Options
		for _, param := range cmd_options {
			str_pmtrs += fmt.Sprintf(" *[%s]*", param.Name)
		}
		str_desc += fmt.Sprintf(
			"%s%s%s\n-# %s\n\n", cmd_prefix, cmd.AppCommand.Name, str_pmtrs, cmd.AppCommand.Description,
		)
	}
	// footer showing our active page
	hcmd_footer := fmt.Sprintf(
		"Page %d of %d",
		page_active,
		page_max,
	)
	embed := &discordgo.MessageEmbed{
		Title:       str_title,
		Description: str_desc,
		Footer: &discordgo.MessageEmbedFooter{
			Text: hcmd_footer,
		},
		Color: 0x41aa0e,
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
			Name:        "help",
			Description: "Show a list and usage of available commands.",
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
