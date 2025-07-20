package terminal

import (
	"bufio"
	"discordgo-bot/utils"
	"discordgo-bot/utils/ucolor"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	Session *discordgo.Session
	fcmds   []TerminalCommand
)

////////////////////////////////////////////////////////////
// Commands ////////////////////////////////////////////////

func init() {
	tcmds := []TerminalCommand{
		{
			Name:        "help",
			Usage:       "[Page]",
			Description: "Show this list.",
			Handle: func(args []string) (bool, error) {
				cmdlimit := 12
				fpage := 1
				fpage_max := int(math.Ceil(
					float64(len(fcmds)) /
						float64(cmdlimit),
				))

				// interpret first arg. as page
				spage := utils.GetSliceStr(args, 0, "1")
				page, err := strconv.Atoi(spage)
				if err == nil && page > 0 {
					fpage = min(page, fpage_max)
				}
				// generate our string to print
				to_print := ucolor.SUBTITLE
				for i, cmd := range fcmds {
					// offset according to our page & visible commands per
					if i < cmdlimit*(fpage-1) {
						continue
					} else if i+1 > cmdlimit*fpage {
						break
					}
					tusage := "" // append usage if we have one
					if len(cmd.Usage) > 1 {
						tusage += " " + ucolor.ITALIC + cmd.Usage + ucolor.RESET + ucolor.SUBTITLE
					}
					to_print += fmt.Sprintf("%s%s - %s\n", cmd.Name, tusage, cmd.Description)
				}

				to_print += fmt.Sprintf("%s%sPage %d of %d%s\n", ucolor.RESET, ucolor.BOLD, fpage, fpage_max, ucolor.RESET)
				fmt.Println(to_print)
				return true, nil
			},
		},
		{
			// fake quit command, actually handled by "interpret" func.
			// we keep this here so our help cmd shows what quit does.
			Name:        "quit",
			Usage:       "",
			Description: "Stop running the bot.",
			Handle:      func(args []string) (bool, error) { return true, nil },
		},
		{
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
		},
		{
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
		}}
	// register all commands listed
	for _, cmd := range tcmds {
		RegisterTerminalCommand(cmd)
	}
}

// Commands ////////////////////////////////////////////////
////////////////////////////////////////////////////////////

type TerminalCommand struct {
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
	signal.Notify(make(chan os.Signal, 1), os.Interrupt)
	fmt.Printf(`
Enter "%shelp%s" for a list of available commands
Quit the program by pressing %sCTRL + D%s or entering "%squit%s".
`,
		ucolor.OKBLUE,
		ucolor.RESET,
		ucolor.OKCYAN,
		ucolor.RESET,
		ucolor.OKBLUE,
		ucolor.RESET,
	)
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
	// ignore whitespace, return on nothing
	if parsed, _ := regexp.MatchString(`\w+`, message); !parsed {
		return 0, nil
	}

	msgSplice := strings.Split(message, " ")
	cmd := strings.ToLower(msgSplice[0])
	args := msgSplice[1:]

	if cmd == "quit" {
		return -1, nil
	}
	for _, tcmd := range fcmds {
		if cmd == tcmd.Name {
			ok, err := tcmd.Handle(args)
			if !ok { // print command usage if formatted wrong
				fmt.Printf("%sUsage: %s %s%s\n", ucolor.BOLD, tcmd.Name, tcmd.Usage, ucolor.RESET)
			}
			return 1, err
		}
	}
	fmt.Printf("%serror: Command \"%s\" not recognized.%s\n", ucolor.FAIL, cmd, ucolor.RESET)
	return 0, nil

}

// Register terminal command.
func RegisterTerminalCommand(cmd TerminalCommand) {
	fcmds = append(fcmds, cmd)
}
