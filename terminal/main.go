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
)

// Start our terminal loop.
func Start() {
	// capture os.Interrupt to prevent hard quitting
	signal.Notify(make(chan os.Signal), os.Interrupt)
	fmt.Println(`
Enter "help" for a list of available commands"
Quit the program with CTRL + D or entering "quit".`)
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
		// handle invalid commands & quitting
		switch code {
		case 100:
			run = false
		case -1:
			fmt.Printf("error: Command \"%s\" not recognized.\n", input)
		}
	}
	fmt.Println("\nQuitting...")

}

// Interpet commands sent via terminal.
//
// Return -1 on failure, 0-1 on success, 100 on quit.
func interpret(message string) (int, error) {

	msgSplice := strings.Split(message, " ")
	cmd := strings.ToLower(msgSplice[0])
	pmtrs := msgSplice[1:]

	// groupcomp := func(s string, str_slice []string) bool {
	// return slices.ContainsFunc(str_slice, func(word string) bool { return s == word })
	// }

	switch cmd {
	case "speak":
		return speak(pmtrs)
	case "quit":
		return 100, nil
	}
	return -1, nil
}

// Send a message in the specified channel.
func speak(parameters []string) (int, error) {
	if len(parameters) < 2 {
		fmt.Println("Usage: speak (channelID) (Message...)")
		return 0, nil
	}
	channel := parameters[0]
	message := strings.Join(parameters[1:], " ")
	_, err := Session.ChannelMessageSend(channel, message)
	return 1, err
}
