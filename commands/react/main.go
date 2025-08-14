// React command
//
// Usage: react [listed response]
package react

import (
	"discordgo-bot/core/commands"
	"discordgo-bot/utils/ucolor"
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	reactions map[string]string = map[string]string{
		"evil":      "https://i.postimg.cc/4dcB8Pqm/baffuievil7.png",
		"angy":      "https://i.postimg.cc/g0nMhmKy/baffuinie7-angy.png",
		"stare":     "https://i.postimg.cc/SRpVH8bV/he-watches.png",
		"bye":       "https://i.postimg.cc/g0CK6Rjd/baffuinie7-by.gif",
		"hello":     "https://i.postimg.cc/sDT94bxT/baffuinie7-hello-fix.gif",
		"jumpscare": "https://i.postimg.cc/DwDrtsqs/ezgif-29dfb41cc8fb1d.gif",
		"pet":       "https://i.postimg.cc/hGdLgjrY/goob.gif",
		"wiggle":    "https://i.postimg.cc/wM6LhkpK/myson.gif",
		"smol":      "https://i.postimg.cc/ZRjCCmg8/baffuinie7-smol.png",
	}
)

func do_command_message(data *commands.DataMessage) {
	var sel string
	parameters := strings.Split(data.Content, " ")
	if len(parameters) > 1 {
		sel = strings.ToLower(parameters[1])
	}

	embed := create_embed(sel)
	_, err := data.Session.ChannelMessageSendComplex(
		data.Message.ChannelID,
		&discordgo.MessageSend{
			Embed: embed,
			Reference: &discordgo.MessageReference{
				MessageID: data.Message.ID,
				ChannelID: data.Message.ChannelID,
				GuildID:   data.Message.GuildID,
			},
			AllowedMentions: &discordgo.MessageAllowedMentions{
				Users:       []string{},
				RepliedUser: false,
			},
		},
	)
	if err != nil {
		log.Panic(err)
	}
}

func do_command_interaction(data *commands.DataInteraction) {
	sel := data.GetOptions()["reaction"].StringValue()

	embed := create_embed(sel)
	err := data.Session.InteractionRespond(data.Interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				embed,
			},
			AllowedMentions: &discordgo.MessageAllowedMentions{
				Users:       []string{},
				RepliedUser: false,
			},
		},
	},
	)
	if err != nil {
		log.Panic(err)
	}
}

func create_embed(selection string) *discordgo.MessageEmbed {
	var embed *discordgo.MessageEmbed
	var img string
	var ok bool

	// get source image from map according to selection
	if img, ok = reactions[selection]; !ok {
		var str_errorsel string
		var str_availablereacts string

		// only show error if a parameter was provided
		if len(selection) > 0 {
			str_errorsel = fmt.Sprintf("-# Reaction \"%s\" does not exist.\n", selection)
		}
		for rname, rimg := range reactions {
			str_availablereacts += fmt.Sprintf("\"[%s](%s)\"\n", rname, rimg)
		}

		embed = &discordgo.MessageEmbed{
			Title: "Baff7 React",
			Description: fmt.Sprintf(`
%s## Available Reactions:
%s`,
				str_errorsel, str_availablereacts),
			Color: ucolor.BAFF_EMBED_DEFAULT,
		}
	} else {
		embed = &discordgo.MessageEmbed{
			Title: " ",
			Image: &discordgo.MessageEmbedImage{
				URL: img,
			},
			Color: ucolor.BAFF_EMBED_DEFAULT,
		}
	}

	return embed
}

func init() {
	// dynamic choice generation
	var command_choices = []*discordgo.ApplicationCommandOptionChoice{}
	for name := range reactions {
		command_choices = append(
			command_choices,
			&discordgo.ApplicationCommandOptionChoice{
				Name:  name,
				Value: name,
			},
		)
	}

	commands.Register(commands.CommandEntry{
		AppCommand: discordgo.ApplicationCommand{
			Name:        "react",
			Description: "Responds with a Baffuinie7 reaction image.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "reaction",
					Description: "Select a reaction.",
					Type:        discordgo.ApplicationCommandOptionString,
					Choices:     command_choices,
					Required:    true,
				},
			},
		},
		FuncMessage:     do_command_message,
		FuncInteraction: do_command_interaction,
	})
}
