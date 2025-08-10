package baff

import (
	"slices"

	"github.com/bwmarrin/discordgo"
)

var (
	BOTID     string   = "510016054391734273"
	CHANNELID []string = []string{"1365436588192043161"}
)

// prevent fatal counting & griefing in the counting channel
func handle_count_read_message(s *discordgo.Session, m *discordgo.MessageCreate) {
	// only read counting bot's messages in the provided channels
	var FORCE bool = true
	if !FORCE && (!slices.Contains(CHANNELID, m.ChannelID) || m.Author.ID != BOTID) {
		return
	}
	// verify we meet the permissions to lock down the channel
	permissions, _ := s.State.UserChannelPermissions(s.State.User.ID, m.ChannelID)
	if permissions&discordgo.PermissionManageChannels == 0 {
		return
	}
}
