package main

import "github.com/bwmarrin/discordgo"

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
