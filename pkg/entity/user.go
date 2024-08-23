package entity

import (
	"github.com/bwmarrin/discordgo"
)

func hasPermission(user discordgo.User) bool {
	member, _ := Session.GuildMember(GuildID, user.ID)
	for _, roleID := range member.Roles {
		if ModeratorID == roleID {
			return true
		}
	}
	return false
}
