// Package provides Discord bot's core functionality.
package core

import (
	"discordgo-bot/baff"
	_ "discordgo-bot/commands"
	"discordgo-bot/core/commands"
	"discordgo-bot/globals"

	"github.com/bwmarrin/discordgo"
)

func Start(session *discordgo.Session) {
	globals.Session = session
	globals.Running = true
	// execute external routines
	commands.InitCommands()
	baff.Start() // baff!
}

func Stop() {
	// shutdown externals before closing session...
	commands.ClearSlashCommands()
	globals.Running = false
}
