package main

import (
	"discordgo-bot-template/commands"
	"discordgo-bot-template/core"
	"discordgo-bot-template/terminal"
	"discordgo-bot-template/utils/ucolor"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"

	_ "discordgo-bot-template/.load"
	_ "discordgo-bot-template/commands"
)

var (
	tokenenv        = "BOT_TOKEN"
	ENABLE_TERMINAL = true
)

func main() {
	// get token from .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Failed to load .env file.\n%s", err)
	}
	token, s := os.LookupEnv(tokenenv)
	if !s {
		log.Fatalf("Token @ %s not found in .env file!", tokenenv)
	}
	// start our discord session
	log.Print("Creating new session...")
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Failed to create new session.\n%s", err)
	}
	session.Open()
	// shutdown processes when exiting
	shutfunc := func() {
		println("Quitting...")
		session.Close()
		core.Stop()
	}
	defer shutfunc()
	// direct to /core
	core.Start(session)
	commands.Ready()

	log.Print("Session successfully launched!")
	// start a terminal cycle
	if ENABLE_TERMINAL && !get_arg("--no-terminal") {
		terminal.Session = session
		terminal.Start()
	} else {
		// capture os.Interrupt to prevent hard quitting
		fmt.Printf(
			"Quit the program by pressing %sCTRL + C%s.\n",
			ucolor.OKCYAN,
			ucolor.RESET,
		)
		a := make(chan os.Signal, 1)
		signal.Notify(a, os.Interrupt)
		<-a
		print("\n")
	}
}

// Check if an arg. matches with provided argument.
func get_arg(arg string) bool {
	for _, a := range os.Args {
		if arg == a {
			return true
		}
	}
	return false
}
