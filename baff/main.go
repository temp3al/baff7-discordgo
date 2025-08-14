package baff

import (
	"discordgo-bot/globals"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	FormalName string = "Baff"
)

func Start() {
	set_status_default()
	//// Discord Call Handlers
	// ./handlers.go
	globals.Session.AddHandler(handle_baff_react)
	globals.Session.AddHandler(handle_gimblo_gang)
	// ./chatbot.go
	globals.Session.AddHandler(handle_readmessage_chatbot)
	// ./countchk.go
	globals.Session.AddHandler(handle_reaction_counting_milestone)
	globals.Session.AddHandler(handle_readmessage_counting_mistake)
}

// Set our activity status.
func set_status_default() {
	var activity []*discordgo.Activity = []*discordgo.Activity{
		{
			Name:    "Baff7 Party",
			Details: "Releasing 2027, probably,,",
			Assets: discordgo.Assets{
				LargeImageID: "https://i.postimg.cc/4NrSMznH/baffled7.jpg",
				SmallImageID: "https://i.postimg.cc/4NrSMznH/baffled7.jpg",
				LargeText:    "Baff7 Party",
				SmallText:    "Baff7 Party",
			},
			Party: discordgo.Party{
				ID:   "245078923450789",
				Size: []int{1, 7},
			},
			Type:      discordgo.ActivityTypeStreaming,
			CreatedAt: time.Now(),
		},
	}

	globals.Session.UpdateStatusComplex(
		discordgo.UpdateStatusData{
			Status:     "online",
			Activities: activity,
		},
	)
}

// Sends a ChatTyping call before executing the provided function.
func TypeBeforeFunction(s *discordgo.Session, channel_id string, time_seconds float64, to_call func()) {
	var time_perping float64 = 6
	var itr int = int(time_seconds / time_perping)
	var time_spillover float64 = time_seconds - (time_perping * float64(itr))
	// typing requests go away after ~9 seconds, so we'll have to send a couple to keep ourselves typing.
	// todo: move this to its own per-channel logic to prevent timeout shenaningans when used in a single channel.
	f_typeping := func(t float64) {
		// time.Duration sucks! To keep the most of our time, convert to Milliseconds
		rt := time.Duration(t*1000) * time.Millisecond
		s.ChannelTyping(channel_id)
		// fmt.Printf("sleeping %v..\n", rt)
		time.Sleep(rt)
	}
	for range itr {
		t := time_perping
		f_typeping(t)
	}
	// wait once more with remaining time before calling
	if time_spillover > 0 {
		f_typeping(time_spillover)
	}
	to_call()
}

// Sends a ChatTyping call until our desired signal is reached.
func TypeBeforeSignal(s *discordgo.Session, cid string) {
	for {
		s.ChannelTyping(cid)
		time.Sleep(3 * time.Second)
	}
}
