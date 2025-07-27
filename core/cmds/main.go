// Package providing slash & chat command functionality.
package cmds

import (
	"discordgo-bot/core"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	ENABLE_CHAT_COMMANDS  bool = true
	ENABLE_SLASH_COMMANDS bool = true
	// Unregister all slash commands out on exit.
	//
	// NOTE: Bots are rate limited to making 200 app commands per day, per guild.
	// Don't use on a big command list unless you need to remove them all, or for cache reasons.
	CLEAR_SLASH_ON_EXIT bool = false
	ChatPrefix               = []string{"+"}

	OptionRegex = []*optRegex{
		{ // opt:"quote clean"
			Regex: regexp.MustCompile(`(\w+):(")([^(")]*)(")`),
			Type:  optParameterCleanType,
		},
		{ // opt:raw_clean
			Regex: regexp.MustCompile(`(\w+):(\w+)`),
			Type:  optParameterCleanType,
		},
		{ // "quote dirty"
			Regex: regexp.MustCompile(`(")([^(")]*)(")`),
			Type:  optParameterDirtyType,
		},
		{ // raw_dirty
			// umm...
			Regex: regexp.MustCompile(`\b(\w+)\b`),
			Type:  optParameterDirtyType,
		},
	}

	command_map = map[string]*CommandEntry{}
	regcmd_map  = map[string]*discordgo.ApplicationCommand{}
)

type CommandEntry struct {
	// 'discordgo.AppCommand' attached to this command.
	//
	// Commands will be automatically generate from it.
	AppCommand discordgo.ApplicationCommand
	// Function to call for this command.
	HandleFunc func(
		createdata *CommandCreateData,
		options map[string]*discordgo.ApplicationCommandInteractionDataOption,
	)
}

type CommandCreateData struct {
	Session     *discordgo.Session
	Message     *discordgo.MessageCreate
	Interaction *discordgo.InteractionCreate
}

// Return our active create variable.
func (s CommandCreateData) GetActive() int {
	if s.Message != nil {
		return CreateMessageType
	} else if s.Interaction != nil {
		return CreateInteractionType
	}
	log.Panicf("neither \"Message\" nor \"Interaction\" are active.")
	return 0
}

const CreateMessageType = 1
const CreateInteractionType = 2

//func (s optParameter) scheck() bool {
//	if s.IndexEndsAt-s.IndexStartsAt == len(s.Content) {
//		return true
//	}
//	return false
//}

func handle_interpret_chatcommand(session *discordgo.Session, message *discordgo.MessageCreate) {
	// ignore bot & app messages (including own)
	if message.Author.Bot {
		return
	}
	var content string
	got_match := false
	// check for chat command prefix
	for _, pref := range ChatPrefix {
		var ok bool
		content, ok = strings.CutPrefix(message.Content, pref)
		if ok {
			got_match = true
			break
		}
	}
	if !got_match {
		return
	}
	// store command for later usage
	command := chcmd_get_commandentry(content)
	options, err := chcmd_get_options(content)
	if err != nil {
		log.Printf("%e\n", err)
		return
	}
	// once everything is ready, execute our command
	command.HandleFunc(
		&CommandCreateData{session, message, nil},
		options,
	)
}

// find and return command from command string
func chcmd_get_commandentry(content string) *CommandEntry {
	if len(content) < 1 {
		log.Fatalf("can't run 'chcmd_match_command' with an empty string")
		return nil
	}
	cmdstr := strings.Split(content, " ")[0] // only interested in the first word from message
	for _, cmd_entry := range command_map {
		if cmdstr == cmd_entry.AppCommand.Name {
			return cmd_entry
		}
	}
	return nil
}

// generate interaction options from command string
func chcmd_get_options(content string) (map[string]*discordgo.ApplicationCommandInteractionDataOption, error) {
	// entries' order is important for priority comparison
	// lower priority regexes will catch words from
	// previous, high priority regexes and are meant to
	// pick up dirtier, less specific option parameters.
	options_map := map[string]*discordgo.ApplicationCommandInteractionDataOption{}

	param_entries := []*optionParameter{}
	// regex, parse & log all parameters unrepeated
	for _, pptr := range OptionRegex {
		rmatch_idx := pptr.Regex.FindAllStringSubmatchIndex(content, -1)
		rmatch_str := pptr.Regex.FindAllStringSubmatch(content, -1)

		for i, rentry := range rmatch_idx {
			idx_start := rentry[0]
			idx_end := rentry[1]
			// check if this match's index range conflicts with a previous one
			// if so, prioritize previously registered option and skip
			do_entry := true
			for _, prm_entry := range param_entries {
				if idx_start <= prm_entry.IndexEndsAt && prm_entry.IndexStartsAt <= idx_end {
					do_entry = false
					break
				}
			}
			if do_entry {
				param_entries = append(param_entries, &optionParameter{
					IndexStartsAt: idx_start,
					IndexEndsAt:   idx_end,
					Content:       rmatch_str[i],
					Type:          pptr.Type,
				})
			}
		}
	}
	// sort low to high
	sort.Slice(param_entries, func(a, b int) bool {
		return param_entries[a].IndexStartsAt < param_entries[b].IndexStartsAt
	})
	// make sure that our clean parameters go after
	// our dirty ones, otherwise, return an error
	is_dirty := true
	for _, parameter := range param_entries {
		if parameter.Type == optParameterCleanType {
			is_dirty = false
		} else if parameter.Type == optParameterDirtyType && !is_dirty {
			return nil, fmt.Errorf("dirty options should always go before clean ones")
		}
	}
	// pack our parameters into proper option arguments
	for _, parameter := range param_entries {
	}
	return options_map, nil
}

type optRegex struct {
	Regex *regexp.Regexp
	Type  int
}

type optionParameter struct {
	IndexStartsAt int
	IndexEndsAt   int
	Content       []string
	Type          int
}

const optParameterCleanType = 1
const optParameterDirtyType = 0

// options' values have to match their assigned
//
//	type or else we'll get struck with a panic error.
func pair_cmdoption_type(v string, t discordgo.ApplicationCommandOptionType) any {
	switch t {
	case discordgo.ApplicationCommandOptionInteger:
		fv, _ := strconv.Atoi(v)
		return float64(fv)
	case discordgo.ApplicationCommandOptionBoolean:
		b, _ := strconv.ParseBool(v)
		return b
	// todo: add handlers to prevent wrongful variables here
	case discordgo.ApplicationCommandOptionUser:
		return v
	case discordgo.ApplicationCommandOptionChannel:
		return v
	case discordgo.ApplicationCommandOptionRole:
		return v
	case discordgo.ApplicationCommandOptionMentionable:
		return v
	case discordgo.ApplicationCommandOptionNumber:
		fv, _ := strconv.ParseFloat(v, 64)
		return float64(fv)
	case discordgo.ApplicationCommandOptionAttachment:
		return nil
	}
	log.Panicf("do not.")
	return nil
}

// option_map := map[string]*discordgo.ApplicationCommandInteractionDataOption{}
//
//	for _, ecmd := range command_map {
//		opt_atlas := ecmd.AppCommand.Options
//		if split_content[0] == ecmd.AppCommand.Name {
//			// attempt to interpret options & convert them to InteractionDataOption
//			for _, opt_i := range split_content[1:] { // skip first (command)
//				for _, option := range opt_atlas {
//					if opt_v, ok := strings.CutPrefix(opt_i, option.Name+":"); ok { // "look for "{opt}:{val}", returns {val}"
//						// if we pass multiple options with the same name, keep
//						// only the first one to prevent any conflicts.
//						// note: should allow for replace or error fallback?
//						if _, ok = option_map[option.Name]; ok {
//							continue
//						}
//						option_map[option.Name] = &discordgo.ApplicationCommandInteractionDataOption{
//							Name:  option.Name,
//							Value: match_value_type(opt_v, option.Type),
//							Type:  option.Type,
//						}
//					}
//				}
//			}
//			// a len zero parmap means we didn't find any labeled
//			// options, attempt parsing them in order
//			if len(option_map) < 1 {
//				for i, opt_i := range split_content[1:] { // skip first (command)
//					if i == len(opt_atlas) {
//						break
//					}
//					option_map[opt_atlas[i].Name] = &discordgo.ApplicationCommandInteractionDataOption{
//						Name:  opt_atlas[i].Name,
//						Value: match_value_type(opt_i, opt_atlas[i].Type),
//						Type:  opt_atlas[i].Type,
//					}
//				}
//			}
//			ecmd.HandleFunc(&CommandCreateData{session, message, nil}, option_map)
//		}
//	}

func handle_interpret_slashcommand(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	// get options and mush them into a map
	opt := interaction.ApplicationCommandData().Options
	opt_map := make(
		map[string]*discordgo.ApplicationCommandInteractionDataOption,
		len(opt))
	for _, option := range opt {
		opt_map[option.Name] = option
	}

	// match a command to our interaction request; run its function handle
	if cmd_entry, ok := command_map[interaction.ApplicationCommandData().Name]; ok {
		cmd_entry.HandleFunc(&CommandCreateData{session, nil, interaction}, opt_map)
	}
}

// Register a new command.
func Register(cmd CommandEntry) {
	cmd_name := cmd.AppCommand.Name
	// only allow registering commands before the bot starts operations
	// it's overall a pretty bad practice to do otherwise
	if core.IsRunning() {
		log.Println("Can't register commands on runtime.")
		return
	}
	// prevent command name conflicts
	if _, ok := command_map[cmd_name]; ok {
		log.Panicf(
			"error: trying to register command with name \"%s\" twice.",
			cmd_name,
		)
	}
	command_map[cmd_name] = &cmd // register our CommandEntry type
}

// Return our list of registered commands.
//
// Useful for help-like commands that require the knowledge of external commands.
func GetCommandEntries() map[string]*CommandEntry {
	return command_map
}

// Start the creation and listening for commands.
func InitCommands() {
	if !ENABLE_CHAT_COMMANDS && !ENABLE_SLASH_COMMANDS {
		return
	}
	if ENABLE_SLASH_COMMANDS {
		log.Println("Registering slash commands...")
		for i, ecmd := range command_map {
			cmd, err := core.CoreSession.ApplicationCommandCreate(core.CoreSession.State.User.ID, "", &ecmd.AppCommand)
			if err != nil {
				log.Panicf("Failed to register command \"%v\": %v", ecmd.AppCommand.Name, err)
			}
			regcmd_map[i] = cmd
		}
		core.CoreSession.AddHandler(handle_interpret_slashcommand)
		log.Print("Slash command registration successful!")
	}
	if ENABLE_CHAT_COMMANDS {
		core.CoreSession.AddHandler(handle_interpret_chatcommand)
		log.Println("Listening for commands in chat!")
	}
}

// Unregister all slash commands out of our bot.
//
// NOTE: Bots are rate limited to making 200 app commands per day, per guild.
// Don't use on a big command list unless you need to remove them all, or for cache reasons.
func SlashClearCommands() {
	if CLEAR_SLASH_ON_EXIT {
		log.Println("Removing commands...")
		for _, v := range regcmd_map {
			err := core.CoreSession.ApplicationCommandDelete(core.CoreSession.State.User.ID, "", v.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}
}
