package bot

import (
	"time"

	"github.com/vatsimnetwork/go-ptd-bot/internal/util"

	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"
)

func AddMemberHandlers(s *discordgo.Session) {
	s.AddHandler(ProcessGuildMemberChunks)
}

func ProcessAllGuilds(s *discordgo.Session) {
	for _, guild := range s.State.Guilds {
		err := s.RequestGuildMembers(guild.ID, "", 0, "1", true)
		if err != nil {
			sentry.CaptureException(err)
		}
	}
}

func ProcessGuildMemberChunks(s *discordgo.Session, mc *discordgo.GuildMembersChunk) {
	for _, member := range mc.Members {
		go util.ProcessMember(s, mc.GuildID, member)
		time.Sleep(45 * time.Millisecond)
	}
}
