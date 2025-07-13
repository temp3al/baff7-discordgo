package tools

import (
	"discordgo-bot/terminal"
	"discordgo-bot/utils"
	"discordgo-bot/utils/ucolor"
	"discordgo-bot/utils/uembed"
	"fmt"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

func init() {
	terminal.TerminalRegcom(terminal.TerminalCom{
		Name:        "test-fail-embed",
		Description: "Test an embed failure message.",
		Usage:       "(ChannelID) [IDNumber]",
		Handle: func(args []string) (bool, error) {
			// cant execute command with less than 1 argument
			if len(args) < 1 {
				return false, nil
			}
			channel := utils.GetSliceStr(args, 0, "X")
			id, _ := strconv.Atoi(utils.GetSliceStr(args, 1, ""))
			_, err := terminal.Session.ChannelMessageSendComplex(
				channel,
				&discordgo.MessageSend{
					Embed: uembed.GenerateErrorMessage(id),
				},
			)
			if err == nil {
				fmt.Printf(
					"%sMessage proccessed @ %s!%s\n",
					ucolor.OKCYAN,
					channel,
					ucolor.RESET,
				)
			}
			return true, err
		},
	})
}
