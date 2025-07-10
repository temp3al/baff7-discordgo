package tools

import (
	"discordgo-bot-template/terminal"
	"discordgo-bot-template/utils/ucolor"
	"discordgo-bot-template/utils/uembed"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func init() {
	terminal.TerminalRegcom(terminal.TerminalCom{
		Name:        "test-fail-embed",
		Description: "Test an embed failure message.",
		Usage:       "(ChannelID) [Reason]",
		Handle: func(args []string) (bool, error) {
			motive := "Unknown"
			if len(args) > 1 {
				motive = strings.Join(args[1:], " ")
			}
			_, err := terminal.Session.ChannelMessageSendComplex(
				args[0],
				&discordgo.MessageSend{
					Embed: uembed.GenerateErrorMessage(motive),
				},
			)
			if err == nil {
				fmt.Printf("%sMessage proccessed @ %s!%s\n", ucolor.OKCYAN, args[0], ucolor.RESET)
			}
			return true, err
		},
	})
}
