package loader

import (
	// Discord Bot commands
	_ "discordgo-bot/commands/help"
	_ "discordgo-bot/commands/magic8"
	_ "discordgo-bot/commands/ping"

	// terminal commands
	_ "discordgo-bot/terminal/tools"
)
