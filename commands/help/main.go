// Help command.
//
// Usage: /help [page]
package help

import (
	"discordgo-bot-template/commands"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func help_command(session *discordgo.Session, message *discordgo.MessageCreate) error {
	cmd_prefix := "/" // commands.Prefixes[0]
	cmds_per_page := 8
	help_page := 1
	max_pages := int(math.Ceil(
		float64(len(commands.Commands)) /
			float64(cmds_per_page),
	))

	// reformat our 2nd parameter as int, set page to parameter
	param := strings.Fields(message.Content)
	if len(param) > 1 {
		page, err := strconv.Atoi(param[1])
		if err == nil && page > 0 {
			help_page = min(page, max_pages)
		}
	}

	hcmd_title := "Command List"
	hcmd_description := ""
	// generate description from commands
	for i, cmd := range commands.Commands {
		if i < cmds_per_page*(help_page-1) { // offset our command list using help_page & commands_per_page
			continue
		} else if i+1 > cmds_per_page*help_page { // stop generating if we go past our cmds limit.
			break
		}
		usage_seg := ""
		if len(cmd.Usage) > 2 {
			usage_seg = " *" + cmd.Usage + "*"
		}
		hcmd_description = hcmd_description +
			fmt.Sprintf("%s%s%s\n-# %s\n\n", cmd_prefix, cmd.Name, usage_seg, cmd.Description)
	}
	hcmd_footer := fmt.Sprintf(
		"Page %d of %d",
		help_page,
		max_pages,
	)

	_, err := session.ChannelMessageSendComplex(
		message.ChannelID,
		&discordgo.MessageSend{
			Embed: &discordgo.MessageEmbed{
				Title:       hcmd_title,
				Description: hcmd_description,
				Footer: &discordgo.MessageEmbedFooter{
					Text: hcmd_footer,
				},
				Color: 0x41aa0e,
			},
			Reference: &discordgo.MessageReference{
				MessageID: message.ID,
				ChannelID: message.ChannelID,
				GuildID:   message.GuildID,
			},
			AllowedMentions: &discordgo.MessageAllowedMentions{
				RepliedUser: false,
			},
		},
	)
	return err
}

func init() {
	commands.Register(
		commands.Command{ // create command
			Name:        "help",
			Description: "Shows you a list of commands.",
			Aliases:     []string{"h"},
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
