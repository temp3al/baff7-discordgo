// 8ball command.
//
// Usage: /8ball (question)
package magic8

import (
	"discordgo-bot/core/cmds"
	"fmt"
	"log"
	"math/rand/v2"

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

func handlecommand(
	cdata *cmds.CommandCreateData,
	parameters map[string]*discordgo.ApplicationCommandInteractionDataOption,
) {
	pmt, ok := parameters["question"]
	if !ok {
		fmt.Println(fmt.Errorf("no question provided?"))
		return
	}
	usr_question := pmt.StringValue()

	// choose a random response from all our answer pools
	res_pool := [][]string{responses_pos, responses_neutro, responses_neg}
	res_group := res_pool[rand.IntN(len(res_pool))]
	response := res_group[rand.IntN(len(res_group))]
	// create an embed including the user's prompt & the outcome

	var author string = "User"
	if cdata.Message != nil {
		author = cdata.Message.Author.GlobalName
	} else if cdata.Interaction != nil {
		author = cdata.Interaction.User.GlobalName
	}

	msg := fmt.Sprintf("**%s's question:** %s\n**answer:** %s", author, usr_question, response)
	embed := &discordgo.MessageEmbed{
		Title:       "8 Ball",
		Description: msg,
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
			Name:        "8ball",
			Description: "Responds to a yes / no question.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "question",
					Description: "What is your question?",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
				},
			},
		},
		HandleFunc: handlecommand,
	})
}
