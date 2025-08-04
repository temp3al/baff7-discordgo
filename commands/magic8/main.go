// 8ball command.
//
// Usage: /8ball (question)
package magic8

import (
	"discordgo-bot/core/commands"
	"fmt"
	"log"
	"math/rand/v2"

	"github.com/bwmarrin/discordgo"
)

// note: could move this to a json file within
// this folder for cleaner customization?
var (
	resp_positive = []string{
		"Sure",
		"Yea",
		"I think so",
		"Thats likely I think",
	}
	resp_neutral = []string{
		"Hmm",
		"Im not quite sure",
		"Eh,, I dont know",
	}
	resp_negative = []string{
		"Probably not",
		"Nooo",
		"I dont think so",
		"I wouldnt say that",
	}
)

func do_command_message(data *commands.DataMessage) {
	question := data.Content
	author := data.Message.Author.DisplayName()

	content := create_content(question, author)
	_, err := data.Session.ChannelMessageSendComplex(
		data.Message.ChannelID,
		&discordgo.MessageSend{
			Content: content,
			Reference: &discordgo.MessageReference{
				MessageID: data.Message.ID,
				ChannelID: data.Message.ChannelID,
				GuildID:   data.Message.GuildID,
			},
			AllowedMentions: &discordgo.MessageAllowedMentions{
				RepliedUser: false,
			},
		},
	)
	if err != nil {
		log.Panic(err)
	}
}

func do_command_interaction(data *commands.DataInteraction) {
	question := data.GetOptions()["question"].StringValue()
	author := data.Interaction.Member.User.DisplayName()

	content := create_content(question, author)
	err := data.Session.InteractionRespond(data.Interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	},
	)
	if err != nil {
		log.Panic(err)
	}
}

func create_content(question string, user string) string {
	var content string

	// pick a response!
	r_ipool := [][]string{resp_positive, resp_neutral, resp_negative}
	r_spool := r_ipool[rand.IntN(len(r_ipool))]
	response := r_spool[rand.IntN(len(r_spool))]
	if rand.IntN(10) >= 7 {
		response += ",,,"
	}
	if rand.IntN(10) >= 3 {
		response += "7"
	}

	content = fmt.Sprintf(
		"-# %s asked: %s\n%s",
		user, question, response,
	)
	return content
}

func init() {
	commands.Register(commands.CommandEntry{
		AppCommand: discordgo.ApplicationCommand{
			Name:        "ask",
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
		FuncMessage:     do_command_message,
		FuncInteraction: do_command_interaction,
	})
}
