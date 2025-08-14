package baff

import (
	"discordgo-bot/utils"
	"math/rand/v2"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	bot_id      string   = "510016054391734273"
	channels_id []string = []string{
		"1365436588192043161",
		"1388969389872906290",
	}
	// ⚠️ @user You have used <n> of your saves. You have <n> left. <motive>
	regex_saveuser regexp.Regexp = *regexp.MustCompile(`⚠️ (.+) You have used (.+) of your saves. You have (.+) left. ?(.+)`)
	// ⚠️ @user You have used <n> guild save! There are <n>/<n> guild saves left. <motive>
	regex_saveguild regexp.Regexp = *regexp.MustCompile(`⚠️ (.+) You have used (.+) guild save! There are (.+)/(.+) guild saves left. ?(.+)`)

	response_milestone []string = []string{
		"<:pog:1405160762816856206><:pog:1405160762816856206><:pog:1405160762816856206> ${COUNT}!!7",
		"omg were at ${COUNT},,",
		"!!!",
		"[${COUNT}!!!!](https://i.postimg.cc/DwDrtsqs/ezgif-29dfb41cc8fb1d.gif)",
		"-# <:pog:1405160762816856206>",
		"Yoo!!7 <a:baffuinie7_wiggle:1365660612436299856>",
	}

	response_usesave_user []string = []string{
		"${USER} please don't do that,,,",
		"${USER} be careful!!",
	}
	response_usesave_guild []string = []string{
		"-# ${BOT} judges ${USER} silently.",
		"${USER} get some saves!7",
		"https://i.postimg.cc/SRpVH8bV/he-watches.png",
	}
)

// React to counting milestones via counting bot reactions.
func handle_reaction_counting_milestone(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if !(r.Member.User.ID == bot_id && slices.Contains([]string{"☑️", "✅", "💯"}, r.Emoji.Name)) {
		return
	}

	m, err := s.ChannelMessage(r.ChannelID, r.MessageID)
	if err != nil {
		return
	}

	msplit := regexp.MustCompile(`(^\d+) ?(.+)?`).FindStringSubmatch(m.Content)
	strnum := msplit[1]
	number, err := strconv.Atoi(strnum)
	// if reacted, shouldnt happen? but
	// retreat to prevent further issues
	if err != nil {
		return
	}
	//var text string
	//if len(msplit) >= 3 {
	//	text = msplit[2]
	//}
	//// please stop bothering me, golang...
	//if len(text) > 0 {
	//}

	reactwith := func(r string) { s.MessageReactionAdd(m.ChannelID, m.ID, r) }

	if chk_multof(number, 1000) {
		// if a multiple of 1000, celebrate!
		reactwith("🎉")
		// get and format a response to use in a sent message in the same channel
		response := utils.RandomStringFromList(response_milestone)
		user := m.Author.Username
		map_replace := map[string]string{
			"${BOT}":   FormalName,
			"${USER}":  user,
			"${COUNT}": msplit[1],
		}
		response = do_format_response(response, map_replace)
		message := &discordgo.MessageSend{
			Content: response,
			Reference: &discordgo.MessageReference{
				MessageID: m.ID,
				ChannelID: m.ChannelID,
				GuildID:   m.GuildID,
			},
		}
		TypeBeforeFunction(s, m.ChannelID, 0.77, func() { s.ChannelMessageSendComplex(m.ChannelID, message) })

	} else if chk_multof(number, 100) {
		// if a multiple of 100, react with :100:
		reactwith("💯")
	} else if slices.Contains([]string{"69", "420", "8008"}, strnum) {
		emojip := []string{
			"<:flushdoggo:1365895504050913371>",
			"<:flabbygasted:1397973649717203096>",
			"<a:catto_boom:1397949570775777401>",
			"<:cowblush:1397723371206082610>",
			"<:squonkler:1397973531811119247>",
			"<:steamhappy:1398377642532671699>",
		}
		reactwith(utils.RandomStringFromList(emojip))
	}
	if strings.Contains(strnum, "7") {
		if rand.Float64() >= 0.7 {
			reactwith("<:baffled:1365660564269039636>")
		}
	}
	reactwith(r.Emoji.Name)
}

func chk_multof(n1 int, n2 int) bool {
	return n1%n2 == 0
}

// React to counting mistakes, prevent fatal counting errors
func handle_readmessage_counting_mistake(s *discordgo.Session, m *discordgo.MessageCreate) {
	// only read counting bot's messages in the provided channels
	var FORCE bool = false
	if !slices.Contains(channels_id, m.ChannelID) || !FORCE && (m.Author.ID != bot_id) || (m.Author.ID == s.State.User.ID) {
		return
	}
	// verify we meet the permissions to lock down the channel
	permissions, _ := s.State.UserChannelPermissions(s.State.User.ID, m.ChannelID)
	if permissions&discordgo.PermissionManageChannels == 0 {
		return
	}

	var match []string

	var user string
	var saves string

	var response string

	// if a user has lost their own save, respond to that with some flavor text.
	if match = regex_saveuser.FindStringSubmatch(m.Content); len(match) > 0 {
		// match: ["⚠️ @user You have used <n> of your saves. You have <n> left. <motive>" "@user" "<n>" "<n>" "<motive>"]
		user = match[1]
		saves = match[2]
		response = utils.RandomStringFromList(response_usesave_user)
	} else if match = regex_saveguild.FindStringSubmatch(m.Content); len(match) > 0 {
		// match: ["⚠️ @user You have used <n> guild save! There are <n>/<n> guild saves left. <motive>" "@user" "<n>" "<n>" "<n>" "<motive>"]
		user = match[1]
		saves = match[2]
		saves_parsed, err := strconv.ParseFloat(saves, 64)
		response = utils.RandomStringFromList(response_usesave_guild)
		// lock down when under 1 guild save
		if err == nil && saves_parsed <= 1 {
			// todo: I don't understand how channel permission overrides work at all, and
			// 2:50am me is NOT in the mood to read through the docs.
		}
	}

	map_replace := map[string]string{
		"${BOT}":   FormalName,
		"${USER}":  user,
		"${SAVES}": saves,
	}

	if len(match) > 0 {
		TypeBeforeFunction(s, m.ChannelID, utils.RandFloat(1.12, 1.6),
			func() { s.ChannelMessageSend(m.ChannelID, do_format_response(response, map_replace)) },
		)
	}

}

func do_format_response(text string, replacement_map map[string]string) string {
	var text_output string = text
	m := replacement_map
	for key, replace := range m {
		text_output = strings.ReplaceAll(text_output, key, replace)
	}
	return text_output
}
