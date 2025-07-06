// Help command.
//
// Usage: /help [page]
package magic8

import (
	"discordgo-bot-template/commands"
	"fmt"
	"math/rand/v2"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// todo: move these to a .json file
var (
	responses_pos = []string{
		"Yes!",
	}
	responses_neutro = []string{
		"Maybe...",
	}
	responses_neg = []string{
		"No.",
	}
)

func magic8_command(session *discordgo.Session, message *discordgo.MessageCreate) error {
	i := strings.Fields(message.Content)
	if len(i) <= 1 {
		return fmt.Errorf("no question")
	}
	// choose a random response from all our answer pools
	res_pool := [][]string{responses_pos, responses_neutro, responses_neg}
	res_group := res_pool[rand.IntN(len(res_pool))]
	response := res_group[rand.IntN(len(res_group))]
	// create an embed including the user's prompt & the outcome
	prompt := strings.Join(i[1:], " ")

	msg := fmt.Sprintf("**%s's question:** %s\n**answer:** %s", message.Author, prompt, response)
	_, err := session.ChannelMessageSendComplex(message.ChannelID, &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Title:       "Magic 8",
			Description: msg,
		},
		Reference: &discordgo.MessageReference{
			MessageID: message.ID,
			ChannelID: message.ChannelID,
			GuildID:   message.GuildID,
		},
		AllowedMentions: &discordgo.MessageAllowedMentions{
			RepliedUser: false,
		},
	})
	return err
}

func init() {
	commands.Register(
		commands.Command{ // create command
			Name:        "8ball",
			Description: "Respond to a yes / no question.",
			Aliases:     []string{"magic8", "answer"},
			// chat message handle
			HandlerChat: func(session *discordgo.Session, message *discordgo.MessageCreate) error {
				return magic8_command(session, message)
			},
			// slash message handle
			HandlerSlash: func(session *discordgo.Session, interaction *discordgo.InteractionCreate) error {
				return nil
			},
		},
	)
}
