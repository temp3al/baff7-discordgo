// Package providing slash & chat command functionality.
package commands

import (
	"discordgo-bot-template/core"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	// Text an user has to prefix to their messages for them to be considered commands
	Prefixes = []string{"+"}

	Commands []Command
)

type Command struct {
	// Name of this command.
	// Can't have spaces, use underscores instead.
	Name string
	// The command's description (Used in slash commands.)
	Description string
	// Parameters the command requires.
	// This is displayed on /help only and doesn't reflect on handler input.
	Usage string
	// Other names that trigger this command.
	// Can't have spaces, use underscores instead.
	Aliases []string
	// Other names to register this slash command under.
	SlashNames []string
	// Function to trigger when calling this command via chat message.
	HandlerChat func(session *discordgo.Session, message *discordgo.MessageCreate) error
	// Function to trigger when calling this command via slash message.
	HandlerSlash func(session *discordgo.Session, interaction *discordgo.InteractionCreate) error
}

func Ready() {
	core.Session.AddHandler(read_command)
}

func read_command(session *discordgo.Session, message *discordgo.MessageCreate) {
	// ignore bot messages (which includes our own)
	if message.Author.Bot {
		return
	}

	// (try to) interpret incoming messages as commands.
	for _, prefix := range Prefixes {
		// store a prefix-less string for later
		msg_no_prefix, match := strings.CutPrefix(message.Content, prefix)

		if match {
			s := strings.Fields(msg_no_prefix)
			if len(s) < 1 { // ignore empty (prefix-only) messages
				continue
			} else if len(s) == 1 { // add empty entry if there
				s = append(s, "") // are no parameters
			}
			// match 1st field with command name or alias
			for _, command := range Commands {
				if strings.EqualFold(s[0], command.Name) {
					err := command.HandlerChat(session, message)
					if err != nil {
						log.Println(err)
					}
				}
			}
		}
	}
}

// Register a new command.
func Register(command Command) {
	Commands = append(Commands, command)
	// register slash command
	_ = []discordgo.ApplicationCommand{
		{
			Name:        command.Name,
			Description: command.Description,
		},
	}
}

// Clear all commands.
func ClearCommands() error {
	return nil
}
