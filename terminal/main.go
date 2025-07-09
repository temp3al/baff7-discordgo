package terminal

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	Session *discordgo.Session
	fcmds   []terminalCmd
)

type terminalCmd struct {
	Name   string
	Handle func(args []string) error
}

// Start our terminal loop.
func Start() {
	// capture os.Interrupt to prevent hard quitting
	signal.Notify(make(chan os.Signal), os.Interrupt)
	fmt.Println(`
Enter "help" for a list of available commands
Quit the program by pressing CTRL + D or entering "quit".`)
	run := true

	clr_in := func(message string) string {
		return strings.TrimSuffix(strings.TrimSuffix(message, "\n"), "\r")
	}

	for {
		if !run {
			break
		}
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil { // os.Interrupt (Ctrl+C) will land us here
			break
		}
		// input will contain a trailing line break as a side
		// effect to our reader: clean it before processing it
		input = clr_in(input)
		code, err := interpret(input)
		if err != nil { // soft print errors
			log.Println(err)
		}
		// handle commands & quitting
		switch code {
		case -1:
			run = false
		}
	}
	print("\n")

}

// Interpet commands sent via terminal.
// Returns (bool, error).
//
//	0: no command match | 1: command match | -1: quit executed
//	error: command failed to execute properly
func interpret(message string) (int, error) {

	msgSplice := strings.Split(message, " ")
	cmd := strings.ToLower(msgSplice[0])
	args := msgSplice[1:]

	for _, tcmd := range fcmds {
		if cmd == tcmd.Name {
			return 1, tcmd.Handle(args)
		}
	}
	if cmd == "quit" {
		return -1, nil
	}
	fmt.Printf("error: Command \"%s\" not recognized.\n", cmd)
	return 0, nil

}

// Register terminal command.
func register_cmd(cmd terminalCmd) {
	fcmds = append(fcmds, cmd)
}

// Commands //

func init() {
	register_cmd(terminalCmd{
		Name: "speak",
		Handle: func(args []string) error {
			if len(args) < 2 {
				fmt.Println("Usage: speak (channelID) (Message...)")
				return nil
			}
			channel := args[0]
			message := strings.Join(args[1:], " ")
			_, err := Session.ChannelMessageSend(channel, message)
			return err
		},
	})
}
