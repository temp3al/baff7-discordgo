package uembed

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var (
	ErrorMessage = discordgo.MessageEmbed{
		Title:       "Error",
		Description: "Something wrong has happened!",
		Color:       0xff1313,
	}
)

func GenerateErrorMessage(motive string) *discordgo.MessageEmbed {
	newErr := &ErrorMessage
	newErr.Description = fmt.Sprintf(`
We're sorry, something went wrong while processing your command...
-# fbstring: %s
-# code: %s
`, motive, "0x0FBBC123")
	return newErr
}
