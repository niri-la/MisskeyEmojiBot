package repository

import (
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

type DiscordRepository interface {
	SendDirectMessage(requestUser string, message string) error
	DeleteChannel(channelID string) error
	CloseThread(threadID string, messageID string) error
	CountMembersWithSpecificRole(guildID string, roleID string) (int, error)
	HasRole(guildID string, user discordgo.User, targetRole string) bool
}

type discordRepository struct {
	session *discordgo.Session
}

func NewDiscordRepository(session *discordgo.Session) DiscordRepository {
	return &discordRepository{session: session}
}

func (r *discordRepository) SendDirectMessage(requestUser string, message string) error {
	user, err := r.session.User(requestUser)
	if err != nil {
		return err
	}
	direct, err := r.session.UserChannelCreate(user.ID)
	if err != nil {
		return err
	}
	r.session.ChannelMessageSend(direct.ID, message)
	return nil
}

func (r *discordRepository) DeleteChannel(channelID string) error {
	_, err := r.session.ChannelDelete(channelID)
	return err
}

func (r *discordRepository) CloseThread(threadID string, messageID string) error {
	channel, _ := r.session.Channel(threadID)
	if !channel.IsThread() {
		return errors.New("channel is not thread")
	}

	archived := true
	locked := true

	t, err := r.session.ChannelEditComplex(channel.ID, &discordgo.ChannelEdit{
		Archived: &archived,
		Locked:   &locked,
	})

	if err != nil {
		return err
	}

	err = r.session.ChannelMessageDelete(t.ParentID, messageID)
	if err != nil {
		return err
	}
	return nil
}

func (r *discordRepository) CountMembersWithSpecificRole(guildID string, roleID string) (int, error) {
	members, err := r.session.GuildMembers(guildID, "", 1000)
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

func (r *discordRepository) HasRole(guildID string, user discordgo.User, targetRole string) bool {
	member, _ := r.session.GuildMember(guildID, user.ID)
	for _, roleID := range member.Roles {
		if targetRole == roleID {
			return true
		}
	}
	return false
}
