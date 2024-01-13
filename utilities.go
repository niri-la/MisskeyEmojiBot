package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

func countMembersWithRole(s *discordgo.Session, guildID string, roleID string) (int, error) {
	members, err := s.GuildMembers(guildID, "", 1000)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, m := range members {
		for _, r := range m.Roles {
			if r == roleID {
				count++
				break
			}
		}
	}

	return count, nil
}

func closeThread(threadID string, messageID string) {
	channel, _ := Session.Channel(threadID)
	if !channel.IsThread() {
		return
	}
	archived := true
	locked := true
	t, err := Session.ChannelEditComplex(channel.ID, &discordgo.ChannelEdit{
		Archived: &archived,
		Locked:   &locked,
	})

	err = Session.ChannelMessageDelete(t.ParentID, messageID)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"event": "delete-thread-message",
		}).Error(err)
	}
}
