package baff

import (
	"discordgo-bot/globals"
	"math/rand/v2"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func Start() {
	globals.Session.AddHandler(handle_baff_react)
	globals.Session.AddHandler(handle_gimblo_gang)
	// custom status
	globals.Session.UpdateStatusComplex(
		discordgo.UpdateStatusData{
			Status: "online",
			Activities: []*discordgo.Activity{
				{
					Name:      "Baff7 Party",
					Details:   "Releasing 2027, probably,,",
					Type:      discordgo.ActivityTypeGame,
					CreatedAt: time.Now(),
				},
			},
		},
	)
	// ./chatbot.go
	globals.Session.AddHandler(handle_chatbot_read_message)
	// ./countchk.go
	globals.Session.AddHandler(handle_count_read_message)
}

// React with a Baffuinie7 related emoji if a sent message contains a Baffuinie7 related keyword.
func handle_baff_react(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot { // ignore bots
		return
	}

	var (
		emojis []string = []string{
			":baffled:1365660564269039636",
			":baffuinie7:1365660581083746346",
			":baffuinie7_angy:1365660594551918675",
			":baffuinie7_baby:1365660715230167162",
			":baffuinie7_blush:1365660758503067700",
			":baffuinie7_chop:1365660775334543441",
			":baffuinie7_confuse:1365660812273909760",
			":baffuinie7_dog:1365660826102403123",
			":baffuinie7_explosion:1365660840115699813",
			":baffuinie7_face:1365660868804739112",
			":baffuinie7_golf:1365660884076331038",
			":baffuinie7_horse:1365660899049865226",
			":baffuinie7_naked:1365660913708830802",
			":baffuinie7_pray:1365660927654891630",
			":baffuinie7_ramjam:1365660946017812540",
			":baffuinie7_sad:1365660964875145216",
			":baffuinie7_scared:1365660978183798835",
			":baffuinie7_sheep:1365660992947753062",
			":baffuinie7_standing:1365661228906840095",
			":baffuinie7_think:1365661251526459413",
			":baffuinie7_think2:1365660628345163786",
			":baffuinie7_victory:1365661280559562833",
			":baffuinie7_wiggle:1365660612436299856",
			":chicken_baffy:1368882908944666714",
		}
		keywords []string = []string{
			"baff",
			"uinie7",
			"affuinie",
		}
	)
	for _, kw := range keywords {
		if strings.Contains(strings.ToLower(m.Content), kw) {
			s.MessageReactionAdd(
				m.ChannelID, m.ID,
				emojis[rand.IntN(len(emojis))],
			)
			return
		}
	}
}

// React to all messages in a specific channel with gimblo emotes.
func handle_gimblo_gang(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot { // ignore bots
		return
	}

	var (
		emojis []string = []string{
			":gimbevil:1365630793254309929",
			":gimbgood:1365630764028268565",
			":gimblo:1355141758681612397",
			":gimblo_cowboy:1367795132526628885",
		}
	)
	var GimbloChannel string = "1355141642939535424"
	if m.ChannelID != GimbloChannel {
		return
	}
	s.MessageReactionAdd(
		m.ChannelID, m.ID,
		emojis[rand.IntN(len(emojis))],
	)
}
