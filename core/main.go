// Package provides Discord bot's core functionality.
package core

import "github.com/bwmarrin/discordgo"

var (
	Session *discordgo.Session
)

func Start(session *discordgo.Session) {
	Session = session
}

func Stop() {
}
