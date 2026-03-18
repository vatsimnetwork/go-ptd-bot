package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"
)

func ProcessAllGuilds(s *discordgo.Session) {
	for _, guild := range s.State.Guilds {
		err := s.RequestGuildMembers(guild.ID, "", 0, "1", true)
		if err != nil {
			sentry.CaptureException(err)
		}
	}
}
