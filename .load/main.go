package loader

import (
	// Discord Bot commands
	_ "discordgo-bot-template/commands/help"
	_ "discordgo-bot-template/commands/magic8"
	_ "discordgo-bot-template/commands/ping"

	// terminal commands
	_ "discordgo-bot-template/terminal/tools"
)
