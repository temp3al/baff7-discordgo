package terminal

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	Session *discordgo.Session
	fcmds   []terminalCmd
)

////////////////////////////////////////////////////////////
// Commands ////////////////////////////////////////////////

func init() {
	register_cmd(terminalCmd{
		Name:        "help",
		Usage:       "",
		Description: "Show this list.",
		Handle: func(args []string) (bool, error) {
			cmdlimit := 8
			fpage := 1
			fpage_max := int(math.Ceil(
				float64(len(fcmds)) /
					float64(cmdlimit),
			))

			if len(args) > 1 {
				page, err := strconv.Atoi(args[1])
				if err == nil && page > 0 {
					fpage = min(page, fpage_max)
				}
			}

			to_print := ""
			for i, cmd := range fcmds {
				if i < cmdlimit*(fpage-1) {
					continue
				} else if i+1 > cmdlimit*fpage {
					break
				}
				tusage := ""
				if len(cmd.Usage) > 1 {
					tusage += fmt.Sprintf(" %s", cmd.Usage)
				}
				to_print += fmt.Sprintf("%s%s - %s\n", cmd.Name, tusage, cmd.Description)
			}
			to_print += fmt.Sprintf("\nPage %d of %d", fpage, fpage_max)
			fmt.Println(to_print)
			return true, nil
		},
	})
	register_cmd(terminalCmd{
		Name:        "speak",
		Usage:       "(ChannelID) (Message...)",
		Description: "Send a message to a channel.",
		Handle: func(args []string) (bool, error) {
			if len(args) < 2 {
				return false, nil
			}
			channel := args[0]
			message := strings.Join(args[1:], " ")
			_, err := Session.ChannelMessageSend(channel, message)
			return true, err
		},
	})
	register_cmd(terminalCmd{
		Name:        "clear",
		Description: "Clear the terminal.",
		Handle: func(args []string) (bool, error) {
			clsFunc := map[string]*exec.Cmd{
				"linux":   exec.Command("clear"),
				"windows": exec.Command("cmd", "/c", "cls"),
			}
			osget := runtime.GOOS
			eCmd, succ := clsFunc[osget]
			if !succ {
				eCmd = clsFunc["linux"]
				fmt.Printf("Your platform \"%s\" is not properly implemented. Attempting fallback...\n", osget)
			}
			eCmd.Stdout = os.Stdout
			return true, eCmd.Run()
		},
	})
}

// Commands ////////////////////////////////////////////////
////////////////////////////////////////////////////////////

type terminalCmd struct {
	Name        string
	Usage       string
	Description string // pref. 1st person
	// returns proper usage and error.
	// if false, handler will print the proper usage of the command.
	Handle func(args []string) (bool, error)
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
			ok, err := tcmd.Handle(args)
			if !ok { // print command usage if formatted wrong
				fmt.Printf("Usage: %s %s\n", tcmd.Name, tcmd.Usage)
			}
			return 1, err
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
